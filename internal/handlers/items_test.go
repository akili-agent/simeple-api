package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jkaninda/okapi"

	"github.com/akili-agent/simple-api/internal/handlers"
	"github.com/akili-agent/simple-api/internal/models"
)

// newItemsApp builds a fresh app with an empty in-memory store per test, so
// tests stay independent of each other.
func newItemsApp() *okapi.Okapi {
	o := okapi.New(okapi.WithOpenAPIDisabled(), okapi.WithAccessLogDisabled())
	g := okapi.NewGroup("/api/v1", o)
	items := handlers.NewItems()
	okapi.RegisterRoutes(o, []okapi.RouteDefinition{
		{Method: http.MethodGet, Path: "/items", Group: g, Handler: items.List, Response: &models.ItemList{}},
		{Method: http.MethodGet, Path: "/items/{id:int}", Group: g, Handler: okapi.H(items.Get), Request: &models.ItemIDRequest{}, Response: &models.Item{}},
		{Method: http.MethodPost, Path: "/items", Group: g, Handler: okapi.H(items.Create), Request: &models.CreateItemRequest{}, Response: &models.Item{}},
	})
	return o
}

// do performs one request against the app and returns the recorder.
func do(app http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var rdr *bytes.Reader
	if body != "" {
		rdr = bytes.NewReader([]byte(body))
	} else {
		rdr = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, rdr)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

func createItem(t *testing.T, app http.Handler, name string) models.Item {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		t.Fatalf("marshaling payload: %v", err)
	}
	rec := do(app, http.MethodPost, "/api/v1/items", string(payload))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d (%s)", rec.Code, rec.Body.String())
	}
	var item models.Item
	if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
		t.Fatalf("create: decoding response: %v", err)
	}
	return item
}

func TestListItemsEmpty(t *testing.T) {
	app := newItemsApp()

	rec := do(app, http.MethodGet, "/api/v1/items", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var list models.ItemList
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if len(list.Items) != 0 {
		t.Fatalf("expected no items, got %d", len(list.Items))
	}
}

func TestCreateItem(t *testing.T) {
	app := newItemsApp()

	item := createItem(t, app, "first item")

	if item.ID != 1 {
		t.Fatalf("expected first item to get ID 1, got %d", item.ID)
	}
	if item.Name != "first item" {
		t.Fatalf("expected name to round-trip, got %q", item.Name)
	}
	if item.CreatedAt.IsZero() {
		t.Fatal("expected the store to set a creation timestamp")
	}
}

func TestCreateItemValidation(t *testing.T) {
	app := newItemsApp()

	tests := []struct {
		name string
		body string
	}{
		{"empty name", `{"name":""}`},
		{"missing name", `{}`},
		{"name too long", `{"name":"` + string(make([]byte, 101)) + `"}`},
		{"malformed json", `{"name":`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := do(app, http.MethodPost, "/api/v1/items", tt.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d (%s)", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	app := newItemsApp()
	created := createItem(t, app, "gadget")

	rec := do(app, http.MethodGet, "/api/v1/items/1", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var item models.Item
	if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if item != created {
		t.Fatalf("expected %+v, got %+v", created, item)
	}
}

func TestGetItemNotFound(t *testing.T) {
	app := newItemsApp()

	rec := do(app, http.MethodGet, "/api/v1/items/99", "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetItemInvalidID(t *testing.T) {
	app := newItemsApp()

	rec := do(app, http.MethodGet, "/api/v1/items/abc", "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestListItemsAfterCreates(t *testing.T) {
	app := newItemsApp()
	createItem(t, app, "one")
	createItem(t, app, "two")

	rec := do(app, http.MethodGet, "/api/v1/items", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var list models.ItemList
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list.Items))
	}
	// The store must answer in ID order so responses are deterministic.
	if list.Items[0].ID != 1 || list.Items[1].ID != 2 {
		t.Fatalf("expected items ordered by ID, got %+v", list.Items)
	}
}
