package main

import (
	"fmt"
	"time"
	"github.com/ifIMust/srsr/client"
)

// Example of how a service should use the provided Go client.
func main() {
	// The client takes care of sending heartbeats, so keep it stored / in scope
	// for the duration of the service's lifetime.
	c := client.NewServiceRegistryClient("flardtech", "http://localhost:4949", "http://localhost:8080")
	err := c.Register()
	if err != nil {
		fmt.Println("Register- error: ", err.Error())
	} else {
		fmt.Println("Registered successfully.")
	}

	// Live a little, let the heartbeats throb
	<- time.After(time.Second * 22)
	
	// A while later, when shutting down...
	c.Deregister()
}
