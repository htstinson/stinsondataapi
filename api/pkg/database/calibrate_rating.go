package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectCalibrateMention(ctx context.Context, subscriber model.Subscriber, search_result_id string) (*[]model.CalibrateMention, error) {
	fmt.Println("d SelectCalibrateMention")

	table := "calibrate_mentions"
	schema_name := subscriber.Schema_Name

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, calibrate_result_id, rating, subscriber_id, title, body,
	rating_date, author, location, headline 
	FROM %s.%s WHERE search_result_id = $1`, schema_name, table)

	rows, err := d.DB.QueryContext(ctx, query, search_result_id)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(rows.Err())

	var items []model.CalibrateMention
	for rows.Next() {
		var item model.CalibrateMention

		err = rows.Scan(&item.ID, &item.CreatedAt, &item.CalibrateResultID, &item.Rating, &item.SubscriberID, &item.Title, &item.Body,
			&item.RatingDate, &item.Author, &item.Location, &item.Headline)

		if err != nil {
			fmt.Println(err.Error())
			return &items, nil
		}

		items = append(items, item)
	}

	return &items, nil
}
