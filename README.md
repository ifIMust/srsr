# srsr
Really Simple Service Registry

## Description
srsr is a service registry written in Go. Its purpose is to allow microservices to find each other.
The primary goals are easy setup, and easy client implementation.

srsr uses [Gin](https://gin-gonic.com/) to offer an HTTP-based API.
It's very trusting; it does little to no validation to prevent abuse.

Clients are provided for Go and Python projects.
By default, clients are expected to send a heartbeat every 30 seconds, or they will be deregistered.
The provided clients are configured to send a heartbeat every 20 seconds.

## Usage

### Server
Precompiled binaries are available for most systems.
```
chmod +x ./srsr-linux-amd64
./srsr-linux-amd64 [-p PORT] [-t TIMEOUT_SECONDS]
```

### Client
A Python client is provided [here](https://github.com/ifIMust/srsrpy).

A Go client is provided in the `client` package.
`go get github.com/ifIMust/srsr/client`
```
import "github.com/ifIMust/srsr/client"
// ...
var c client.ServiceRegistry = client.NewServiceRegistry(server_address, my_name, my_address)
c.Register()
// ...
c.Deregister()
```

## API Endpoints
All actions are performed as JSON Post requests.

### /register
Register a service. The client should do this once at startup, and store the returned ID for later use.
Example request:
```
{"name": "flard_service", "address": "321.123.321.123:4321"}
```
Response:
```
{"success": true, "id": "1ccda9cb-0432-4306-965d-6e0fbad571bc"}
```

### /deregister
Deregister a service. The client should do this once at shutdown, using the ID stored from registration.
Example request:
```
{"id": "1ccda9cb-0432-4306-965d-6e0fbad571bc"}
```
(Response body unused)

### /lookup
Retrieve an address for a service.
Example request:
```
{"name": "flard_service"}
```
Response examples:
```
{"success": "true", "address": "321.123.321.123:4321"}

{"success": "false", "address": ""}
```

### /heartbeat
Inform the service registry that the client is still up, to avoid automatic deregistration.
Example request:
```
{"id": "1ccda9cb-0432-4306-965d-6e0fbad571bc"}
```
(Response body unused)


## Further plans
- Improve API tests for 100% coverage.
- Test suite for Go client.
