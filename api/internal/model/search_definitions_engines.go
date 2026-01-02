package model

import "time"

type SearchDefinitionEngines struct {
	Id                  string    `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	ModifiedAt          time.Time `json:"modified_at"`
	SearchEngineId      string    `json:"search_engine_id"`
	SearchDefinitionsId string    `json:"search_definitions_id"`
}

type SearchDefinitionEnginesView struct {
	Id                   string    `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	ModifiedAt           time.Time `json:"modified_at"`
	SearchEngineId       string    `json:"search_engine_id"`
	SearchDefinitionsId  string    `json:"search_definitions_id"`
	SearchEngineName     string    `json:"search_engine_name"`
	SearchDefinitionName string    `json:"search_definition_name"`
	SearchQuery          string    `json:"search_query"`
	EngineId             string    `json:"engine_id"`
	DefinitionId         string    `json:"definition_id"`
}
