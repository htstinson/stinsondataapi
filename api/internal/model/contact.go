package model

import "time"

// Contact
type Contact struct {
	Id             string    `json:"id"`
	ParentId       string    `json:"parent_id"` // The parent_id is the customer's id value.
	LastName       string    `json:"last_name"`
	FirstName      string    `json:"first_name"`
	Schema_Name_   string    `json:"schema_name"`
	Subscriber_Id_ string    `json:"subscriber_id"`
	CreatedAt      time.Time `json:"created_at"`
	ModifiedAt     time.Time `json:"modified_at"`
}
