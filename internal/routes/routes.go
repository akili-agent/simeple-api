// Package routes declares every HTTP route as data ([]okapi.RouteDefinition)
// so the routing table doubles as the OpenAPI definition.
package routes

import (
	"net/http"

	"github.com/jkaninda/okapi"

	"github.com/akili-agent/simple-api/internal/handlers"
	"github.com/akili-agent/simple-api/internal/models"
)

// basePath is the versioned prefix every API route lives under.
const basePath = "/api/v1"

// Routes returns the full routing table of the service. The Okapi instance is
// only used to build the /api/v1 group; registration stays in the caller so
// tests can inspect or register the same table on their own instance.
func Routes(o *okapi.Okapi) []okapi.RouteDefinition {
	api := okapi.NewGroup(basePath, o)
	items := handlers.NewItems()

	return []okapi.RouteDefinition{
		{
			Method:      http.MethodGet,
			Path:        "/healthz",
			Group:       api,
			Handler:     handlers.Healthz,
			OperationId: "get-healthz",
			Summary:     "Liveness probe",
			Description: "Reports whether the service process is alive.",
			Tags:        []string{"probes"},
			Response:    &models.Status{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/readyz",
			Group:       api,
			Handler:     handlers.Readyz,
			OperationId: "get-readyz",
			Summary:     "Readiness probe",
			Description: "Reports whether the service is ready to receive traffic.",
			Tags:        []string{"probes"},
			Response:    &models.Status{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/items",
			Group:       api,
			Handler:     items.List,
			OperationId: "list-items",
			Summary:     "List items",
			Description: "Returns every item currently stored, ordered by ID.",
			Tags:        []string{"items"},
			Response:    &models.ItemList{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/items/{id:int}",
			Group:       api,
			Handler:     okapi.H(items.Get),
			OperationId: "get-item",
			Summary:     "Get an item",
			Description: "Returns the item with the given ID, or 404 when it does not exist.",
			Tags:        []string{"items"},
			Request:     &models.ItemIDRequest{},
			Response:    &models.Item{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/items",
			Group:       api,
			Handler:     okapi.H(items.Create),
			OperationId: "create-item",
			Summary:     "Create an item",
			Description: "Creates a new item and answers 201 Created with it.",
			Tags:        []string{"items"},
			Request:     &models.CreateItemRequest{},
			Response:    &models.Item{},
		},
	}
}
