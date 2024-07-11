# srsr
Really Simple Service Registry

## Description
HTTP-based registry service. Post json to register and lookup other services.
It's really simple and very trusting; it does little to no validation to prevent abuse.

This project uses [Gin](https://gin-gonic.com/).

## Endpoints
### /register (POST)
Register a service. Example:
```
{"name": "flard_service", "address": "321.123.321.123:4321"}
```

Response:
```
{"success": true, id="1ccda9cb-0432-4306-965d-6e0fbad571bc"}
```

### /lookup (POST)
Retrieve an address for a service. Example:
```
{"name": "flard_service"}
```

Response:
```
{"success": "true", "address": "321.123.321.123:4321"}
{"success": "false"}
```
