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
	SubscriberID             uuid.UUID  `json:"subscriber_id" db:"subscriber_id"`
	Published                *time.Time `json:"published"`
}

type CalibrateSearchResultView struct {
	ResultId                 uuid.UUID  `json:"id"`
	Link                     *string    `json:"link,omitempty"`
	Snippet                  *string    `json:"snippet,omitempty"`
	Title                    *string    `json:"title,omitempty"`
	SearchTime               *time.Time `json:"search_time,omitempty"`
	ResultCreatedAt          time.Time  `json:"created_at"`
	SubscriberId             uuid.UUID  `json:"subscriber_id"`
	SearchDefinitionId       uuid.UUID  `json:"search_definition_id"`
	SearchDefinitionName     *string    `json:"search_definition_name"`
	Query                    *string    `json:"query"`
	SearchDefinitionComment  *string    `json:"search_definition_comment"`
	ExactMatch               bool       `json:"exact_match"`
	MaxResults               int        `json:"max_results"`
	SortByDate               bool       `json:"sort_by_date"`
	StartDate                *time.Time `json:"start_date"`
	EndDate                  *time.Time `json:"end_date"`
	SearchType               *string    `json:"search_type"`
	SearchEngineId           uuid.UUID  `json:"search_engine_id"`
	SearchEngineName         *string    `json:"search_engine_name"`
	SearchEngineIdentifier   *string    `json:"search_engine_identifier"`
	SearchEngineComment      *string    `json:"search_engine_comment"`
	SearchDefinitionEngineID *uuid.UUID `json:"search_definition_engine_id,omitempty" db:"search_definition_engine_id"`
	Published                *time.Time `json:"published"`
}
