package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) ListSearchEngines(ctx context.Context, subscriber model.Subscriber, limit, offset int) ([]model.SearchEngine, error) {

	fmt.Println("d ListSearchEngines")

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, name, search_engine_Id, comment FROM %s.calibrate_search_engines 
	WHERE parent_id = '%s' ORDER BY name ASC LIMIT $1 OFFSET $2`, subscriber.Schema_Name, subscriber.Id)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing search_engines: %w", err)
	}
	defer rows.Close()

	var searchengines []model.SearchEngine
	for rows.Next() {
		var searchengine model.SearchEngine
		if err := rows.Scan(&searchengine.Id, &searchengine.CreatedAt, &searchengine.ModifiedAt, &searchengine.Name, &searchengine.SearchEngineId, &searchengine.Comment); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning search_definition: %w", err)
		}

		searchengines = append(searchengines, searchengine)
	}
	return searchengines, nil
}
