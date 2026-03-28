package model

import "time"

// Customer
type Profile struct {
	Id         string    `json:"id"`
	ParentId   string    `json:"parentid"`
	Legal_Name string    `json:"legal_name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
