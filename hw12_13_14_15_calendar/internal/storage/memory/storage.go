package memorystorage

import (
	"fmt"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"sync"
)

type Storage struct {
	items map[string]*storage.Event
	mu    sync.RWMutex
}

var instance *Storage = nil
var items = map[string]*storage.Event{}

func New() *Storage {
	if instance == nil {
		items = make(map[string]*storage.Event)
		instance = &Storage{}
		instance.items = items
	}

	return instance
}

func (s *Storage) Add(e *storage.Event) error {
	id := uuid.New()
	e.ID = id
	s.items[id.String()] = e

	return nil
}

func (s *Storage) Get(id string) (*storage.Event, error) {
	s.mu.RLock()
	item, ok := items[id]
	s.mu.RUnlock()

	if ok == false {
		return nil, fmt.Errorf("error occured: %w", &storage.EventNotFoundErr{Id: id})
	}

	return item, nil
}

func (s *Storage) Delete(id string) {
	delete(s.items, id)
}

func (s *Storage) All() map[string]*storage.Event {
	return s.items
}

func (s *Storage) Edit(id string, e *storage.Event) error {
	s.mu.RLock()
	_, ok := items[id]
	s.mu.RUnlock()

	if ok == false {
		return fmt.Errorf("error occured: %w", &storage.EventNotFoundErr{Id: id})
	}

	s.mu.Lock()
	items[id] = e
	s.mu.Unlock()

	return nil
}
