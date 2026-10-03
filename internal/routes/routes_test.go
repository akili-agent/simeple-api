package routes_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jkaninda/okapi"

	"github.com/akili-agent/simple-api/internal/routes"
)

// newApp registers the production routing table on a fresh Okapi instance.
func newApp() *okapi.Okapi {
	o := okapi.New(okapi.WithAccessLogDisabled())
	okapi.RegisterRoutes(o, routes.Routes(o))
	o.WithOpenAPIDocs(okapi.OpenAPI{Title: "simple-api", Version: "test", UI: okapi.ScalarUI})
	return o
}

// TestRoutesRegisteredUnderAPIV1 makes sure the declared routing table is
// actually served under /api/v1.
func TestRoutesRegisteredUnderAPIV1(t *testing.T) {
	app := newApp()

	for _, path := range []string{
		"/api/v1/healthz",
		"/api/v1/readyz",
		"/api/v1/items",
	} {
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s: expected 200, got %d", path, rec.Code)
		}
	}
}

// TestOpenAPIDocumentsRoutes ensures the Request/Response types declared on
// the routes flow into the generated specification.
func TestOpenAPIDocumentsRoutes(t *testing.T) {
	app := newApp()

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	spec, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatalf("reading spec: %v", err)
	}
	for _, fragment := range []string{
		"/api/v1/healthz",
		"/api/v1/readyz",
		"/api/v1/items",
		"/api/v1/items/{id}",
		"Liveness probe",
		"List items",
		"Create an item",
	} {
		if !strings.Contains(string(spec), fragment) {
			t.Errorf("OpenAPI spec is missing %q", fragment)
		}
	}
}

// TestScalarUIServed ensures the interactive documentation is available.
func TestScalarUIServed(t *testing.T) {
	app := newApp()

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/docs", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatalf("reading docs page: %v", err)
	}
	if !strings.Contains(string(body), "scalar") {
		t.Error("expected /docs to render the Scalar UI")
	}
}
