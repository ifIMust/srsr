package main

import (
	"flag"
	"strconv"
	"time"
	"github.com/ifIMust/srsr/registry"
	"github.com/ifIMust/srsr/server"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 4214, "The server will listen on this port.")
	var timeoutSeconds int
	flag.IntVar(&timeoutSeconds, "t", 30, "Heartbeat timeout (seconds). Clients will be deregistered after this period, if they don't send a heartbeat.")
	flag.Parse()
	registry := registry.NewServiceRegistry()
	registry.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	router := server.SetupRouter(registry)
	router.Run("localhost:" + strconv.Itoa(port))
}
