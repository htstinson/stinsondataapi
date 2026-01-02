package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectSearchDefinitionEnginesView(ctx context.Context, subscriber model.Subscriber, limit, offset int) ([]model.SearchDefinitionEnginesView, error) {
	fmt.Println("d SelectSearchDefinitionEnginesView")

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, search_engine_Id, search_engine_name, 
	search_definition_name, search_query, engine_id, definition_id FROM %s.search_definition_engines_view ORDER BY name ASC LIMIT $1 OFFSET $2`, subscriber.Schema_Name)

	fmt.Println(query)

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
			row.SearchDefinitionName, row.SearchQuery, row.EngineId, row.DefinitionId); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning search_definition: %w", err)
		}

		results = append(results, row)
	}

	return results, nil
}
