package message

type RegisterRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
	Port    string `json:"port"`
}

type RegisterResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type DeregisterRequest struct {
	ID string `json:"id" binding:"required"`
}

type DeregisterResponse struct {
	Success bool `json:"success"`
}

type LookupRequest struct {
	Name string `json:"name" binding:"required"`
}

type LookupResponse struct {
	Success bool   `json:"success"`
	Address string `json:"address"`
}

type HeartbeatRequest struct {
	ID string `json:"id" binding:"required"`
}

type HeartbeatResponse struct {
	Success bool `json:"success"`
}
