// Package store provides the in-memory persistence behind the items resource.
package store

import (
	"sort"
	"sync"
	"time"

	"github.com/akili-agent/simple-api/internal/models"
)

// ItemStore keeps items in process memory. It is safe for concurrent use and
// is meant as a stand-in for a real database in this scaffold.
type ItemStore struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]models.Item
}

// NewItemStore returns an empty store.
func NewItemStore() *ItemStore {
	return &ItemStore{nextID: 1, items: make(map[int64]models.Item)}
}

// List returns every item ordered by ID so responses are deterministic.
func (s *ItemStore) List() []models.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]models.Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

// Get returns the item with the given ID, or false when it does not exist.
func (s *ItemStore) Get(id int64) (models.Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	return item, ok
}

// Create stores a new item with a server-assigned ID and creation timestamp.
func (s *ItemStore) Create(name string) models.Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := models.Item{
		ID:        s.nextID,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
	s.items[item.ID] = item
	s.nextID++
	return item
}
