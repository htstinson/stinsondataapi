package model

type RDSLogin struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	DdbInstanceIdentifier string `json:"dbInstanceIdentifier"`
}

// All
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

type Roles struct {
	Id       string
	Username string
	Names    string
}
