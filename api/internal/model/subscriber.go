package model

import "time"

// Subscriber
type Subscriber struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	Schema_Name string    `json:"schema_name"`
}
