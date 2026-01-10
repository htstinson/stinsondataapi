package model

import (
	"time"

	"github.com/google/uuid"
)

type CalibrateSearchResult struct {
	ID                       uuid.UUID  `json:"id" db:"id"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt               time.Time  `json:"modified_at" db:"modified_at"`
	Link                     *string    `json:"link,omitempty" db:"link"`
	Snippet                  *string    `json:"snippet,omitempty" db:"snippet"`
	Title                    *string    `json:"title,omitempty" db:"title"`
	SearchDefinitionEngineID *uuid.UUID `json:"search_definition_engine_id,omitempty" db:"search_definition_engine_id"`
	SearchTime               *time.Time `json:"search_time,omitempty" db:"search_time"`
}
