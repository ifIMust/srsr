package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ifIMust/srsr/registry"
)

type RegisterRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type RegisterResponse struct {
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

type deregister_request struct {
	ID string `json:"id"`
}

type deregister_response struct {
	Success bool `json:"success"`
}

func register(c *gin.Context, sr registry.Registry) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, reg_err := sr.Register(request.Name, request.Address)
	if reg_err != nil {
		c.AbortWithError(http.StatusBadRequest, reg_err)
	}

	r := RegisterResponse{ID: id, Success: true}
	c.JSON(http.StatusOK, r)
}

func deregister(c *gin.Context, sr registry.Registry) {
	var request deregister_request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reg_err := sr.Deregister(request.ID)
	r := deregister_response{}
	if reg_err == nil {
		r.Success = true
	}
	c.JSON(http.StatusOK, r)
}

func lookup(c *gin.Context, sr registry.Registry) {
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

func SetupRouter(registry registry.Registry) *gin.Engine {
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		register(c, registry)
	})
	router.POST("/deregister", func(c *gin.Context) {
		deregister(c, registry)
	})
	router.POST("/lookup", func(c *gin.Context) {
		lookup(c, registry)
	})
	return router
}
