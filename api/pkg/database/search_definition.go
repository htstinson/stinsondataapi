package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectSearchDefinitions(ctx context.Context, subscriber model.Subscriber, limit, offset int) ([]model.SearchDefinition, error) {

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
