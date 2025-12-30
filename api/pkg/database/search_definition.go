package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectSearchDefinitions(ctx context.Context, customer model.Customer, limit, offset int) ([]model.SearchDefinition, error) {

	fmt.Println("d SelectSearchDefinitions")

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, search_name, search_description, search_definition FROM %s.calibrate_search_definition 
	WHERE parent_id = '%s' ORDER BY lastname,firstname ASC LIMIT $1 OFFSET $2`, customer.Schema_Name, customer.Id)

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
		if err := rows.Scan(&searchdefinition.Id, &searchdefinition.CreatedAt, &searchdefinition.ModifiedAt, &searchdefinition.SearchName, &searchdefinition.SearchDescription, &searchdefinition.SearchDefinition); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning search_definition: %w", err)
		}

		searchdefinitions = append(searchdefinitions, searchdefinition)
	}
	return searchdefinitions, nil
}
