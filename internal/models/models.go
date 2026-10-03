// Package models holds the request and response types shared by handlers and
// routes. Keeping them together lets both sides reference the same types
// without an import cycle.
package models

import "time"

// Status is the body returned by the health and readiness probes.
type Status struct {
	Status string `json:"status" example:"ok" description:"Service health status"`
}

// Item is the example resource served by this API.
type Item struct {
	ID        int64     `json:"id" description:"Unique item identifier"`
	Name      string    `json:"name" description:"Human readable item name"`
	CreatedAt time.Time `json:"createdAt" description:"Creation timestamp (UTC)"`
}

// ItemList wraps the items collection so the response stays a JSON object,
// which leaves room for pagination metadata later without breaking clients.
type ItemList struct {
	Items []Item `json:"items"`
}

// ItemIDRequest carries the {id} path parameter for single-item routes.
type ItemIDRequest struct {
	ID int64 `json:"id" path:"id" required:"true" description:"Item identifier"`
}

// CreateItemRequest carries the body for POST /items.
type CreateItemRequest struct {
	Body struct {
		Name string `json:"name" required:"true" minLength:"1" maxLength:"100" description:"Name of the item to create" example:"My first item"`
	}
}
