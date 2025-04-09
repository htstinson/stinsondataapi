package model

import (
	"time"
)

type RDSLogin struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	DdbInstanceIdentifier string `json:"dbInstanceIdentifier"`
}

// Item
type Blocked struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// Item
type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// User
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never send password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
	Roles        string    `json:"roles"`
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
