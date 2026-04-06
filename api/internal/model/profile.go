package model

import "time"

// Customer
type Profile struct {
	Id             string    `json:"id"`
	ParentId       string    `json:"parentid"`
	Legal_Name     *string   `json:"legal_name"`
	Phone          *string   `json:"phone"`
	Fax            *string   `json:"fax"`
	Email          *string   `json:"email"`
	Website        *string   `json:"website"`
	LinkedIn       *string   `json:"linkedin"`
	Facebook       *string   `json:"facebook"`
	Instagram      *string   `json:"instagram"`
	X              *string   `json:"x"`
	YouTube        *string   `json:"youtube"`
	Pinterest      *string   `json:"pinterest"`
	GoogleBusiness *string   `json:"google_business"`
	Yelp           *string   `json:"yelp"`
	GlassDoor      *string   `json:"glassdoor"`
	Github         *string   `json:"github"`
	NextDoor       *string   `json:"nextdoor"`
	CreatedAt      time.Time `json:"created_at"`
	ModifiedAt     time.Time `json:"modified_at"`
}
