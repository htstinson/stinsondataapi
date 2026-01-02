package model

import "time"

type SearchDefinition struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Name       string    `json:"name"`
	Query      string    `json:"query"`
	Comment    string    `json:"common"`
	ExactMatch bool      `json:"exact_match"`
	MaxResults int       `json:"max_results"`
	SortByDate bool      `json:"sort_by_date"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	SearchType string    `json:"search_type"`
}
