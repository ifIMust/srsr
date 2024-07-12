package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"github.com/ifIMust/srsr/message"
)

const contentType = "application/json"

// Default service timeout is expected to be 30 seconds.
const heartbeatInterval = 20 * time.Second

type ServiceRegistryClient interface {
	Register() error
	Deregister()
}

type client struct {
	serverAddress string
	
	clientName string
	clientAddress string
	
	clientID string
	isRegistered bool
	cancel chan int
}

func NewServiceRegistryClient(clientName string, clientAddress string, serverAddress string) ServiceRegistryClient {
	return &client{
		serverAddress: serverAddress,
		clientName: clientName,
		clientAddress: clientAddress,
		cancel: make(chan int),
	}
}

func (c *client) sendHeartbeat() {
	request := message.HeartbeatRequest{
		ID: c.clientID,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(request)
	if err != nil {
		return
	}
	
	resp, err := http.Post(c.serverAddress + "/heartbeat", contentType, buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()	
}

func (c *client) Register() error {
	if c.isRegistered {
		return errors.New("Register- already registered!")
	}
	
	request := message.RegisterRequest{
		Name: c.clientName,
		Address: c.clientAddress,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.serverAddress + "/register", contentType, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Register- bad status: " + resp.Status)
	}

	response := message.RegisterResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	c.clientID = response.ID
	c.isRegistered = true

	go func() {
		for c.isRegistered {
			select {
			case <-time.After(heartbeatInterval):
				c.sendHeartbeat()
			case <-c.cancel:
				c.isRegistered = false
			}
		}
	}()
	
	return err
}

func (c *client) Deregister() {
	c.cancel <- 1
	
	request := message.DeregisterRequest{
		ID: c.clientID,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(request)
	if err != nil {
		return
	}
	
	resp, err := http.Post(c.serverAddress + "/deregister", contentType, buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
