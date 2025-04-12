package model

import "time"

// Blocked
type Blocked struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}
