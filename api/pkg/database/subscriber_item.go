package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectSubscriberItemView(ctx context.Context, subscriber_id string, limit int, offset int) ([]model.Subscriber_Item_View, error) {
	fmt.Println("database.go SelectSubscriberItem()")

	where_clause := " "

	if subscriber_id != "" {
		_, err := ValidateUUID(subscriber_id)
		if err != nil {
			fmt.Printf("Invalid UUID error: %v\n", err)
			return nil, err
		} else {
			fmt.Println(subscriber_id, "validated.")
			where_clause = fmt.Sprintf(` where subscriber_id = '%s' `, subscriber_id)
		}
	}

	fmt.Println("new query")

	query := fmt.Sprintf("SELECT id, item_id, subscriber_id, item_name, subscriber_name FROM subscriber_items_view%sORDER BY subscriber_name, item_name ASC LIMIT $1 OFFSET $2", where_clause)

	rows, err := d.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var subscriber_item_views []model.Subscriber_Item_View
	for rows.Next() {
		var subscriber_item_view model.Subscriber_Item_View
		if err := rows.Scan(&subscriber_item_view.Id, &subscriber_item_view.Id, &subscriber_item_view.Subscriber_Id, &subscriber_item_view.Item_Name, &subscriber_item_view.Subscriber_Name); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_subscriber: %w", err)
		}

		subscriber_item_views = append(subscriber_item_views, subscriber_item_view)

	}
	return subscriber_item_views, nil
}
