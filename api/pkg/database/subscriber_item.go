package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
		if err := rows.Scan(&subscriber_item_view.Id, &subscriber_item_view.Item_ID, &subscriber_item_view.Subscriber_Id, &subscriber_item_view.Item_Name, &subscriber_item_view.Subscriber_Name); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_subscriber: %w", err)
		}

		fmt.Println(subscriber_item_view.Id)

		subscriber_item_views = append(subscriber_item_views, subscriber_item_view)

	}
	return subscriber_item_views, nil
}

func (d *Database) CreateSubscriberItem(ctx context.Context, item_id string, subscriber_id string) (*model.Subscriber_Item, error) {
	fmt.Println("d CreateSubscriberItem")

	subscriber_item := &model.Subscriber_Item{
		Id:            uuid.New().String(),
		Item_ID:       item_id,
		Subscriber_Id: subscriber_id,
	}

	query := `
        INSERT INTO subscriber_items (id, item_id, subscriber_id) VALUES ($1, $2, $3)
    `

	_, err := d.DB.ExecContext(ctx, query,
		subscriber_item.Id,
		subscriber_item.Item_ID,
		subscriber_item.Subscriber_Id,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating subscriber_item: %w", err)
	}

	return subscriber_item, nil

}

func (d *Database) LookupSubscriberItem(ctx context.Context, item_id string, subscriber_id string) (*model.Subscriber_Item, error) {
	fmt.Println("d LookupUserSubscriber")

	var subscriber_item = model.Subscriber_Item{}

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, item_id, subscriber_id FROM subscriber_items WHERE item_id = $1 and subscriber_id = $2",
		item_id, subscriber_id,
	).Scan(&subscriber_item.Id, &subscriber_item.Item_ID, &subscriber_item.Subscriber_Id)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting subscriber_item: %w", err)
	}

	return &subscriber_item, nil
}

func (d *Database) DeleteSubscriberItem(ctx context.Context, id string) error {
	fmt.Println("d DeleteSubscriberItem")
	fmt.Println(id)

	query := `DELETE FROM subscriber_items WHERE id = $1`

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}

func (d *Database) GetSubscriberItem(ctx context.Context, id string) (*model.Subscriber_Item, error) {
	fmt.Println("d GetSubscriberItem")
	fmt.Println(id)

	var subscriberitem model.Subscriber_Item

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, item_id, subscriber_id FROM subscriber_items WHERE id = $1",
		id,
	).Scan(&subscriberitem.Id, &subscriberitem.Item_ID, &subscriberitem.Subscriber_Id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting subscriber_item: %w", err)
	}

	fmt.Println(subscriberitem.Id, subscriberitem.Item_ID, subscriberitem.Subscriber_Id)

	return &subscriberitem, nil
}
