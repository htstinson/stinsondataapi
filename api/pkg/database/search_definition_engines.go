package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectSearchDefinitionEnginesView(ctx context.Context, search_definition model.SearchDefinition, limit, offset int) ([]model.SearchDefinitionEnginesView, error) {
	fmt.Println("d SelectSearchDefinitionEnginesView")

	subscriber, err := d.GetSubscriber(ctx, search_definition.SubscriberId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, search_engine_Id, search_engine_name, search_definition_name, 
		search_query, engine_id, definition_id FROM %s.search_definition_engines_view WHERE definition_id = '%s' 
		ORDER BY search_engine_name ASC LIMIT $1 OFFSET $2`, subscriber.Schema_Name, search_definition.Id)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing search_definition_engines_view: %w", err)
	}
	defer rows.Close()

	var results []model.SearchDefinitionEnginesView
	for rows.Next() {
		var row model.SearchDefinitionEnginesView
		if err := rows.Scan(&row.Id, &row.CreatedAt, &row.ModifiedAt, &row.SearchEngineId, &row.SearchEngineName,
			&row.SearchDefinitionName, &row.SearchQuery, &row.EngineId, &row.DefinitionId); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning search_definition: %w", err)
		}

		results = append(results, row)
	}

	return results, nil
}

func (d *Database) SelectSearchDefinitionEnginesSubscriberView(ctx context.Context, subscriber model.Subscriber, limit, offset int) ([]model.SearchDefinitionEnginesView, error) {
	fmt.Println("d SelectSearchDefinitionEnginesSubscriberView")

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, search_engine_Id, search_engine_name, search_definition_name, 
		search_query, engine_id, definition_id FROM %s.search_definition_engines_view 
		ORDER BY search_engine_name ASC LIMIT $1 OFFSET $2`, subscriber.Schema_Name)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing search_definition_engines_view: %w", err)
	}
	defer rows.Close()

	var results []model.SearchDefinitionEnginesView
	for rows.Next() {
		var row model.SearchDefinitionEnginesView
		if err := rows.Scan(&row.Id, &row.CreatedAt, &row.ModifiedAt, &row.SearchEngineId, &row.SearchEngineName,
			&row.SearchDefinitionName, &row.SearchQuery, &row.EngineId, &row.DefinitionId); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning search_definition: %w", err)
		}

		results = append(results, row)
	}

	return results, nil
}

func (d *Database) GetSearchDefinitionEnginesView(ctx context.Context, subscriber model.Subscriber, search_definitions_engines_id string) (model.SearchDefinitionEnginesView, error) {
	fmt.Println("d GetSearchDefinitionEnginesView")

	limit := 1
	offset := 0

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, search_engine_Id, search_engine_name, search_definition_name, 
		search_query, engine_id, definition_id FROM %s.search_definition_engines_view WHERE id = $1
		ORDER BY search_engine_name ASC LIMIT $2 OFFSET $3`, subscriber.Schema_Name)

	var row model.SearchDefinitionEnginesView

	rows, err := d.DB.QueryContext(ctx, query, search_definitions_engines_id, limit, offset)
	if err != nil {
		fmt.Println(err.Error())
		return row, fmt.Errorf("error selecting search_engine: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row model.SearchDefinitionEnginesView
		if err := rows.Scan(&row.Id, &row.CreatedAt, &row.ModifiedAt, &row.SearchEngineId, &row.SearchEngineName,
			&row.SearchDefinitionName, &row.SearchQuery, &row.EngineId, &row.DefinitionId); err != nil {
			fmt.Println(err.Error())
			return row, fmt.Errorf("error scanning search_definition_engines_view: %w", err)
		}

	}
	return row, nil
}

func (d *Database) DeleteSearchDefinitionEngine(ctx context.Context, subscriber *model.Subscriber, id string) error {

	fmt.Println("d DeleteSearchDefinitionEngine")

	query := fmt.Sprintf(`DELETE FROM %s.search_definition_engines WHERE id = $1`, subscriber.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

func (d *Database) CreateSearchDefinitionEngine(ctx context.Context, subscriber model.Subscriber, row model.SearchDefinitionEngines) (*model.SearchDefinitionEngines, error) {
	fmt.Println("d CreateSearchDefinitionEngine")

	table := "search_definition_engines"
	schema_name := subscriber.Schema_Name

	fmt.Println(row.SearchEngineId)
	fmt.Println(row.SearchDefinitionsId)

	query := fmt.Sprintf(`INSERT INTO %s.%s (id, search_engine_id, search_definitions_id) 
		VALUES ($1, $2, $3)`, schema_name, table)

	_, err := d.DB.ExecContext(ctx, query,
		row.Id, row.SearchEngineId, row.SearchDefinitionsId)

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating search definition engine: %w", err)
	}

	return &row, nil
}
