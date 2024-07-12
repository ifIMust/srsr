package message

type RegisterRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type RegisterResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type DeregisterRequest struct {
	ID string `json:"id"`
}

type DeregisterResponse struct {
	Success bool `json:"success"`
}

type LookupRequest struct {
	Name string `json:"name"`
}

type LookupResponse struct {
	Success bool   `json:"success"`
	Address string `json:"address"`
}

type HeartbeatRequest struct {
	ID string `json:"id"`
}

type HeartbeatResponse struct {
	Success bool `json:"success"`
}
