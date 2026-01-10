package database

import (
	"context"
	"fmt"
	"time"

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

func (d *Database) SelectSearchResultView(ctx context.Context, subscriber model.Subscriber, sde string) (*[]model.CalibrateSearchResultView, error) {
	fmt.Println("d SelectSearchResultView")

	table := "v_calibrate_search_results"
	schema_name := subscriber.Schema_Name

	query := fmt.Sprintf(`SELECT result_id, link, snippet, title, search_time, result_created_at, subscriber_id,
	search_definition_id, search_definition_name, query, search_definition_comment, exact_match, max_results, sort_by_date,
	start_date, end_date, search_type, search_engine_id, search_engine_name, search_engine_identifier, search_engine_comment,
	search_definition_engine_id, published FROM %s.%s WHERE search_definition_engine_id = $1`, schema_name, table)

	rows, err := d.DB.QueryContext(ctx, query, sde)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(rows.Err())

	var items []model.CalibrateSearchResultView
	for rows.Next() {
		var item model.CalibrateSearchResultView
		if err := rows.Scan(&item.ResultId, &item.Link, &item.Snippet, &item.Title, &item.SearchTime, &item.ResultCreatedAt, &item.SubscriberId,
			&item.SearchDefinitionId, &item.SearchDefinitionName, &item.Query, &item.SearchDefinitionComment, &item.ExactMatch, &item.MaxResults, &item.SortByDate,
			&item.StartDate, &item.EndDate, &item.SearchType, &item.SearchEngineId, &item.SearchEngineName, &item.SearchEngineIdentifier, &item.SearchEngineComment,
			&item.SearchDefinitionEngineID, &item.Published); err != nil {
			return nil, fmt.Errorf("error scanning item: %w", err)
		}
		fmt.Println("published", item.Published.Format(time.RFC3339))
		items = append(items, item)
	}
	return &items, nil
}
