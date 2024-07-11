package main

import (
	"math/rand"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type service_entry struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type register_request struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type register_response struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type lookup_request struct {
	Name string `json:"name"`
}

type lookup_response struct {
	Success bool   `json:"success"`
	Address string `json:"address"`
}

type service_registry struct {
	mutex     sync.Mutex
	Store     map[string]*service_entry
	nameStore map[string][]*service_entry
}

func NewServiceRegistry() *service_registry {
	sr := service_registry{}
	// map ID to entry
	sr.Store = make(map[string]*service_entry)
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

	s.Store[id] = &entry
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

func register(c *gin.Context, sr *service_registry) {
	var request register_request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, reg_err := sr.Register(request.Name, request.Address)
	if reg_err != nil {
		// TODO possibly a different HTTP code in this case
		c.AbortWithError(http.StatusBadRequest, reg_err)
	}

	r := register_response{ID: id, Success: true}
	c.JSON(http.StatusOK, r)
}

func lookup(c *gin.Context, sr *service_registry) {
	var request lookup_request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := sr.Lookup(request.Name)
	r := lookup_response{}
	if len(address) > 0 {
		r.Success = true
		r.Address = address
	}
	c.JSON(http.StatusOK, r)
}

func main() {
	registry := NewServiceRegistry()

	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		register(c, registry)
	})

	router.POST("/lookup", func(c *gin.Context) {
		lookup(c, registry)
	})

	router.Run("localhost:8080")
}
