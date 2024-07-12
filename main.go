package main

import (
	"github.com/ifIMust/srsr/registry"
	"github.com/ifIMust/srsr/server"
)

func main() {
	registry := registry.NewServiceRegistry()
	router := server.SetupRouter(registry)
	router.Run("localhost:8080")
}
