package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectSearchDefinitions(ctx context.Context, subscriber model.Subscriber, limit int, offset int) ([]model.SearchDefinition, error) {

	fmt.Println("d SelectSearchDefinitions")

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, name, comment, query, exact_match, max_results, sort_by_date, start_date, end_date, 
	search_type, subscriber_id FROM %s.calibrate_search_definition 
	ORDER BY name ASC LIMIT $1 OFFSET $2`, subscriber.Schema_Name)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing search_definitions: %w", err)
	}
	defer rows.Close()

	var searchdefinitions []model.SearchDefinition
	for rows.Next() {
		var searchdefinition model.SearchDefinition
		if err := rows.Scan(&searchdefinition.Id, &searchdefinition.CreatedAt, &searchdefinition.ModifiedAt, &searchdefinition.Name, &searchdefinition.Comment,
			&searchdefinition.Query, &searchdefinition.ExactMatch, &searchdefinition.MaxResults, &searchdefinition.SortByDate, &searchdefinition.StartDate,
			&searchdefinition.EndDate, &searchdefinition.SearchType, &searchdefinition.SubscriberId); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning search_definition: %w", err)
		}

		searchdefinitions = append(searchdefinitions, searchdefinition)
	}
	return searchdefinitions, nil
}

func (d *Database) GetSearchDefinition(ctx context.Context, subscriber model.Subscriber, definition_id string, limit int, offset int) (model.SearchDefinition, error) {

	fmt.Println("d GetSearchDefinition")

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, name, comment, query, exact_match, max_results, sort_by_date, start_date, end_date, 
	search_type, subscriber_id FROM %s.calibrate_search_definition WHERE id='%s'
	ORDER BY name ASC LIMIT $1 OFFSET $2`, subscriber.Schema_Name, definition_id)

	var searchdefinition model.SearchDefinition

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return searchdefinition, fmt.Errorf("error listing search_definitions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(&searchdefinition.Id, &searchdefinition.CreatedAt, &searchdefinition.ModifiedAt, &searchdefinition.Name, &searchdefinition.Comment,
			&searchdefinition.Query, &searchdefinition.ExactMatch, &searchdefinition.MaxResults, &searchdefinition.SortByDate, &searchdefinition.StartDate,
			&searchdefinition.EndDate, &searchdefinition.SearchType, &searchdefinition.SubscriberId); err != nil {
			fmt.Println(err.Error())
			return searchdefinition, fmt.Errorf("error scanning search_definition: %w", err)
		}

	}
	return searchdefinition, nil
}

func (d *Database) DeleteSearchDefinition(ctx context.Context, subscriber *model.Subscriber, search_definition_id string) error {

	fmt.Println("d DeleteSearchDefinition")

	query := fmt.Sprintf(`DELETE FROM %s.calibrate_search_definition WHERE id = $1`, subscriber.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query, search_definition_id)
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

func (d *Database) CreateSearchDefinition(ctx context.Context, subscriber model.Subscriber, row model.SearchDefinition) (*model.SearchDefinition, error) {
	fmt.Println("d CreateSearchDefinition")

	table := "calibrate_search_definition"
	schema_name := subscriber.Schema_Name

	query := fmt.Sprintf(`INSERT INTO %s.%s (id, name, comment, query, exact_match, max_results, sort_by_date, start_date, end_date, search_type, subscriber_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, schema_name, table)

	if row.SearchType == "" {
		row.SearchType = "custom"
	}

	_, err := d.DB.ExecContext(ctx, query,
		row.Id, row.Name, row.Comment, row.Query, row.ExactMatch, row.MaxResults, row.SortByDate, row.StartDate, row.EndDate, row.SearchType, row.SubscriberId)

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating search definition: %w", err)
	}

	return &row, nil
}

func (d *Database) UpdateSearchDefinition(ctx context.Context, subscriber *model.Subscriber, row model.SearchDefinition) (*model.SearchDefinition, error) {
	fmt.Println("d UpdateSearchDefinition")

	table := "calibrate_search_definition"
	schema_name := subscriber.Schema_Name

	query := fmt.Sprintf(`UPDATE %s.%s SET name = $1 WHERE id = $2`, schema_name, table)

	fmt.Println(query)

	_, err := d.DB.ExecContext(ctx, query, row.Name, row.Id)

	return &row, err
}
