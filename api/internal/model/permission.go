package model

import "time"

// Permission
type Permission struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Object_Id   string    `json:"object_id"`
	CreatedAt   time.Time `json:"created_at"`
}
