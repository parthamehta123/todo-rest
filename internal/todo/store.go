package todo

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	List() []Item
	Get(id int64) (Item, error)
	Create(title string) Item
	Update(id int64, title *string, done *bool) (Item, error)
	Delete(id int64) error
}

type MemoryStore struct {
	mu    sync.RWMutex
	items map[int64]Item
	next  int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: make(map[int64]Item), next: 1}
}

func (s *MemoryStore) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Item, 0, len(s.items))
	for _, it := range s.items {
		out = append(out, it)
	}
	return out
}

func (s *MemoryStore) Get(id int64) (Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	if !ok {
		return Item{}, ErrNotFound
	}
	return it, nil
}

func (s *MemoryStore) Create(title string) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.next
	s.next++
	it := Item{ID: id, Title: title, CreatedAt: time.Now()}
	s.items[id] = it
	return it
}

func (s *MemoryStore) Update(id int64, title *string, done *bool) (Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[id]
	if !ok {
		return Item{}, ErrNotFound
	}
	if title != nil {
		it.Title = *title
	}
	if done != nil {
		it.Done = *done
	}
	s.items[id] = it
	return it, nil
}

func (s *MemoryStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
