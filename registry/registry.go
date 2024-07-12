package registry

import (
	"errors"
	"math/rand"
	"sync"
	
	"github.com/google/uuid"
)

type Registry interface {
	Register(name string, address string) (string, error)
	Deregister(id string) error
	Lookup(name string) string
}

type service_entry struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type service_registry struct {
	mutex     sync.Mutex
	store     map[string]*service_entry
	nameStore map[string][]*service_entry
}

func NewServiceRegistry() *service_registry {
	sr := service_registry{}
	// map ID to entry
	sr.store = make(map[string]*service_entry)
	// map name to entries
	sr.nameStore = make(map[string][]*service_entry)
	return &sr
}

func (s *service_registry) Register(name string, address string) (string, error) {
	id := uuid.NewString()
	entry := service_entry{ID: id,
		Name:    name,
		Address: address}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[id] = &entry
	_, ok := s.nameStore[name]
	if !ok {
		s.nameStore[name] = make([]*service_entry, 0, 1)
	}
	s.nameStore[name] = append(s.nameStore[name], &entry)
	return id, nil
}

func (s *service_registry) Lookup(name string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	entries, ok := s.nameStore[name]
	if ok {
		return entries[rand.Intn(len(entries))].Address
	}
	return ""
}

func (s *service_registry) Deregister(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	idEntry, ok := s.store[id]
	if ok {
		delete(s.store, id)
		delete(s.nameStore, idEntry.Name)
		return nil
	}
	return errors.New("Deregister - no match for ID")
}
