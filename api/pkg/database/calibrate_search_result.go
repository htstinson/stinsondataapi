package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) CreateSearchResult(ctx context.Context, subscriber model.Subscriber, row model.CalibrateSearchResult) (*model.CalibrateSearchResult, error) {
	fmt.Println("d CreateSearchResult")

	table := "calibrate_search_results"
	schema_name := subscriber.Schema_Name

	query := fmt.Sprintf(`INSERT INTO %s.%s (id, link, snippet, title, search_definition_engine_id, search_time, subscriber_id ) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, schema_name, table)

	id := uuid.New().String()

	_, err := d.DB.ExecContext(ctx, query,
		id, row.Link, row.Snippet, row.Title, row.SearchDefinitionEngineID, row.SearchTime, row.SubscriberID)

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating search result: %w", err)
	}

	return &row, nil
}
