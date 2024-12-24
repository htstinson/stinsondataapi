package model

import "time"

type Test struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never send password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

type AccountQueryResponse struct {
	TotalSize int       `json:"totalSize"`
	Done      bool      `json:"done"`
	Records   []Account `json:"records"`
}

type Account struct {
	Attributes  AccountAttributes `json:"attributes"`
	Id          string            `json:"Id"`
	Name        string            `json:"Name"`
	Industry    *string           `json:"Industry"` // Using pointer since it can be null
	Description string            `json:"description"`
	Phone       string            `json:"phone"`
}

type AccountAttributes struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
