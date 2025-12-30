package model

import "time"

type SearchDefinition struct {
	Id                string    `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"`
	SearchName        string    `json:"search_name"`
	SearchDefinition  string    `json:"search_definition"`
	SearchDescription string    `json:"search_description"`
}
