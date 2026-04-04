package model

import "time"

// Customer
type Address struct {
	Id           string    `json:"id"`
	SubscriberId string    `json:"subscriber_id"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
	AddressType  string    `json:"address_type"`
	AddressUse   string    `json:"address_use"`
	Street1      string    `json:"street1"`
	Street2      string    `json:"street2"`
	POBox        string    `json:"po_box"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Zip          string    `json:"zip"`
}
