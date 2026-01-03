package model

import "time"

type SearchEngine struct {
	Id             string    `json:"id"`
	SubscriberId   string    `json:"subscriber_id"`
	CreatedAt      time.Time `json:"created_at"`
	ModifiedAt     time.Time `json:"modified_at"`
	Name           string    `json:"name"`
	SearchEngineId string    `json:"search_engine_id"`
	Comment        *string   `json:"comment"`
}
