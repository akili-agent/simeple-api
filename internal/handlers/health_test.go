package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jkaninda/okapi"

	"github.com/akili-agent/simple-api/internal/handlers"
	"github.com/akili-agent/simple-api/internal/models"
)

// newProbeApp exposes the probe handlers on their production paths.
func newProbeApp() *okapi.Okapi {
	o := okapi.New(okapi.WithOpenAPIDisabled(), okapi.WithAccessLogDisabled())
	g := okapi.NewGroup("/api/v1", o)
	okapi.RegisterRoutes(o, []okapi.RouteDefinition{
		{Method: http.MethodGet, Path: "/healthz", Group: g, Handler: handlers.Healthz, Response: &models.Status{}},
		{Method: http.MethodGet, Path: "/readyz", Group: g, Handler: handlers.Readyz, Response: &models.Status{}},
	})
	return o
}

func TestHealthz(t *testing.T) {
	app := newProbeApp()

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body models.Status
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("expected status ok, got %q", body.Status)
	}
}

func TestReadyz(t *testing.T) {
	app := newProbeApp()

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/readyz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body models.Status
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("expected status ok, got %q", body.Status)
	}
}
