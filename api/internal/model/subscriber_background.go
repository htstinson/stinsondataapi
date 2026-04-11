package model

import "time"

type Background struct {
	Id           string    `json:"id"`
	SubscriberId string    `json:"subscriber_id"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
	Topic        *string   `json:"topic"`
	Summary      *string   `json:"summary"`
	Details      *string   `json:"details"`
}
