// Package handlers contains the HTTP handlers. Handlers stay thin: binding
// and validation are done by okapi via the request types, and all persistence
// goes through the store.
package handlers

import (
	"fmt"

	"github.com/jkaninda/okapi"

	"github.com/akili-agent/simple-api/internal/models"
	"github.com/akili-agent/simple-api/internal/store"
)

// Healthz answers the liveness probe. It must stay dependency-free so an
// orchestrator can always tell the process is alive.
func Healthz(c *okapi.Context) error {
	return c.OK(models.Status{Status: "ok"})
}

// Readyz answers the readiness probe. There is no external dependency to
// check in this scaffold, so readiness mirrors liveness.
func Readyz(c *okapi.Context) error {
	return c.OK(models.Status{Status: "ok"})
}

// Items groups the handlers of the items example resource over an
// ItemStore, so routes can be declared without package-level state.
type Items struct {
	store *store.ItemStore
}

// NewItems returns the items handlers backed by a fresh in-memory store.
func NewItems() *Items {
	return &Items{store: store.NewItemStore()}
}

// List returns every stored item.
func (h *Items) List(c *okapi.Context) error {
	return c.OK(models.ItemList{Items: h.store.List()})
}

// Get returns a single item or 404 when the ID is unknown. Binding of the
// {id} path parameter is done by okapi.H from models.ItemIDRequest.
func (h *Items) Get(c *okapi.Context, in *models.ItemIDRequest) error {
	item, ok := h.store.Get(in.ID)
	if !ok {
		return c.AbortNotFound(fmt.Sprintf("item %d not found", in.ID))
	}
	return c.OK(item)
}

// Create validates the request body (via tags on models.CreateItemRequest),
// stores the item and answers 201 Created.
func (h *Items) Create(c *okapi.Context, in *models.CreateItemRequest) error {
	return c.Created(h.store.Create(in.Body.Name))
}
