package model

import "time"

// Customer
type Customer struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Parent_Id     string    `json:"parent_id"`
	Subscriber_Id string    `json:"subscriber_id"`
	Schema_Name   string    `json:"schema_name"`
	CreatedAt     time.Time `json:"created_at"`
}
