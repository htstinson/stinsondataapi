package model

import (
	"myapi/internal/salesforce"
	"time"
)

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

// Salesforce
type AccountQueryResponse struct {
	TotalSize int       `json:"totalSize"`
	Done      bool      `json:"done"`
	Records   []Account `json:"records"`
}

type Account struct {
	Attributes       AccountAttributes         `json:"attributes"`
	Id               string                    `json:"Id"`
	Name             string                    `json:"Name"`
	Industry         *string                   `json:"Industry"` // Using pointer since it can be null
	Description      string                    `json:"Description"`
	Phone            string                    `json:"Phone"`
	LastModifiedDate salesforce.SalesforceTime `json:"LastModifiedDate"`
}

type AccountAttributes struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
