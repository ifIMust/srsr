package registry

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

const defaultTimeout = 30 * time.Second

type Registry interface {
	Register(name string, address string) (string, error)
	Deregister(id string) error
	Lookup(name string) string
	Heartbeat(id string)
	SetTimeout(duration time.Duration)
}

type service_entry struct {
	ID      string
	Name    string
	Address string
	Cancel  chan int
	Reset   chan int
}

type service_registry struct {
	mutex          sync.Mutex
	store          map[string]*service_entry
	nameStore      map[string][]*service_entry
	serviceTimeout time.Duration
}

func NewServiceRegistry() *service_registry {
	sr := service_registry{}
	// map ID to entry
	sr.store = make(map[string]*service_entry)
	// map name to entries
	sr.nameStore = make(map[string][]*service_entry)
	sr.serviceTimeout = defaultTimeout
	return &sr
}

func (s *service_registry) Register(name string, address string) (string, error) {
	id := uuid.NewString()
	entry := service_entry{ID: id,
		Name:    name,
		Address: address,
		Cancel:  make(chan int),
		Reset:   make(chan int),
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[id] = &entry
	_, ok := s.nameStore[name]
	if !ok {
		s.nameStore[name] = make([]*service_entry, 0, 1)
	}
	s.nameStore[name] = append(s.nameStore[name], &entry)

	// Goroutine to handle automatic deregistration, in the absence of heartbeats
	go func() {
		var keepGoing = true
		timer := time.NewTimer(s.serviceTimeout)

		for keepGoing {
			select {
			case <-timer.C:
				s.Deregister(id)

			case <-entry.Reset:
				timer.Reset(s.serviceTimeout)

			case <-entry.Cancel:
				keepGoing = false
				timer.Stop()
			}
		}
	}()

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

func (s *service_registry) Heartbeat(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	entry, ok := s.store[id]
	if ok {
		entry.Reset <- 1
	}
}

func (s *service_registry) SetTimeout(duration time.Duration) {
	s.serviceTimeout = duration
}
