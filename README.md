# srsr
A Really Simple Service Registry

## Description
srsr is a service registry written in Go. The primary goals are easy setup, and easy client implementation.
srsr uses [Gin](https://gin-gonic.com/) to offer a fast HTTP-based API.
It's very trusting; it does little to no validation to prevent abuse.

By default, clients are expected to send a heartbeat every 30 seconds, or they will be deregistered.

## Endpoints
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
Coming soon.
Example request:
```
{"id": "1ccda9cb-0432-4306-965d-6e0fbad571bc"}
```
(Response body unused)
