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

type DeregisterRequest struct {
	ID string `json:"id"`
}

type DeregisterResponse struct {
	Success bool `json:"success"`
}

func deregister(c *gin.Context, sr registry.Registry) {
	var request DeregisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reg_err := sr.Deregister(request.ID)
	r := DeregisterResponse{}
	if reg_err == nil {
		r.Success = true
	}
	c.JSON(http.StatusOK, r)
}

type LookupRequest struct {
	Name string `json:"name"`
}

type LookupResponse struct {
	Success bool   `json:"success"`
	Address string `json:"address"`
}

func lookup(c *gin.Context, sr registry.Registry) {
	var request LookupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := sr.Lookup(request.Name)
	r := LookupResponse{}
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
