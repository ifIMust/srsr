package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ifIMust/srsr/message"
	"github.com/ifIMust/srsr/registry"
)

const defaultScheme = "http://"

func register(c *gin.Context, sr registry.Registry) {
	var request message.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// This might be a wrong guess, since it will not include the client port
	if request.Address == "" {
		request.Address = defaultScheme + c.ClientIP()
	}

	id, reg_err := sr.Register(request.Name, request.Address)
	if reg_err != nil {
		c.AbortWithError(http.StatusBadRequest, reg_err)
		return
	}

	r := message.RegisterResponse{ID: id, Success: true}
	c.JSON(http.StatusOK, r)
}

func deregister(c *gin.Context, sr registry.Registry) {
	var request message.DeregisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reg_err := sr.Deregister(request.ID)
	r := message.DeregisterResponse{}
	if reg_err == nil {
		r.Success = true
	}
	c.JSON(http.StatusOK, r)
}

func lookup(c *gin.Context, sr registry.Registry) {
	var request message.LookupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := sr.Lookup(request.Name)
	r := message.LookupResponse{}
	if len(address) > 0 {
		r.Success = true
		r.Address = address
	}
	c.JSON(http.StatusOK, r)
}

func heartbeat(c *gin.Context, sr registry.Registry) {
	var request message.HeartbeatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := message.HeartbeatResponse{Success: sr.Heartbeat(request.ID)}
	c.JSON(http.StatusOK, r)
}

func SetupRouter(registry registry.Registry) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

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
	router.POST("/heartbeat", func(c *gin.Context) {
		heartbeat(c, registry)
	})
	return router
}
