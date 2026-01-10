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

	query := fmt.Sprintf(`INSERT INTO %s.%s (id, ) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, schema_name, table)

	id := uuid.New().String()

	_, err := d.DB.ExecContext(ctx, query,
		row.Id, row.Name, row.Comment, row.Query, row.ExactMatch, row.MaxResults, row.SortByDate, row.StartDate, row.EndDate, row.SearchType, row.SubscriberId)

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating search definition: %w", err)
	}

	return &row, nil
}
