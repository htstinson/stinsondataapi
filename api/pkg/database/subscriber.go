package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// subscribers
func (d *Database) GetSubscriber(ctx context.Context, id string) (*model.Subscriber, error) {
	fmt.Println("d GetSubscriber")

	var subscriber model.Subscriber

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, name, created_at FROM subscribers WHERE id = $1",
		id,
	).Scan(&subscriber.Id, &subscriber.Name, &subscriber.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting subscriber: %w", err)
	}

	return &subscriber, nil
}

func (d *Database) GetSubscriberByName(ctx context.Context, name string) (*model.Subscriber, error) {
	fmt.Printf("[%v] [GetSubscriberByName] %s\n", time.Now().Format(time.RFC3339), name)

	subscriber := &model.Subscriber{}
	query := `
        SELECT id, name, created_at FROM subscribers WHERE username = $1
    `

	err := d.DB.QueryRowContext(ctx, query, name).Scan(
		&subscriber.Id,
		&subscriber.Name,
		&subscriber.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting subscriber: %w", err)
	}

	return subscriber, nil
}

func (d *Database) CreateSubscriber(ctx context.Context, name string) (*model.Subscriber, error) {
	fmt.Println("d CreateSubscriber")

	subscriber := &model.Subscriber{
		Id:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO subscribers (id, name, created_at) VALUES ($1, $2, $3)
    `

	_, err := d.DB.ExecContext(ctx, query,
		subscriber.Id,
		subscriber.Name,
		subscriber.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating subscriber: %w", err)
	}

	return subscriber, nil

}

func (d *Database) SelectSubscribers(ctx context.Context, limit, offset int) ([]model.Subscriber, error) {

	fmt.Println("database.go SelectSubscribers()")

	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, name, created_at, schema_name FROM subscribers ORDER BY name ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing subscribers: %w", err)
	}
	defer rows.Close()

	var subscribers []model.Subscriber
	for rows.Next() {
		var subscriber model.Subscriber
		if err := rows.Scan(&subscriber.Id, &subscriber.Name, &subscriber.CreatedAt, &subscriber.Schema_Name); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning subscriber: %w", err)
		}

		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}

func (d *Database) UpdateSubscriber(ctx context.Context, subscriber *model.Subscriber) error {
	fmt.Println("d UpdateSubscriber")

	query := `UPDATE subscribers SET name = $1 WHERE id = $2`

	_, err := d.DB.ExecContext(ctx, query, subscriber.Name, subscriber.Id)

	return err
}

func (d *Database) DeleteSubscriber(ctx context.Context, subscriber *model.Subscriber) error {
	fmt.Println("d DeleteSubscriber")
	fmt.Println(subscriber.Id)

	fmt.Println("delete user_subscribe_role")
	query := `DELETE from common.user_subscriber_role where user_subscriber_id in
				(select id from common.user_subscriber where subscriber_id = $1);`
	result, err := d.DB.ExecContext(ctx, query, subscriber.Id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result.RowsAffected())
	}

	fmt.Println("delete user subscriber")
	query = `DELETE from common.user_subscriber where subscriber_id = $1;`
	result, err = d.DB.ExecContext(ctx, query, subscriber.Id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result.RowsAffected())
	}

	fmt.Println("delete subscriber")
	query = `DELETE FROM common.subscribers WHERE id = $1`
	result, err = d.DB.ExecContext(ctx, query, subscriber.Id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result.RowsAffected())
	}

	fmt.Println("drop schema")
	query = fmt.Sprintf(`DROP SCHEMA %s cascade`, subscriber.Schema_Name)
	fmt.Println(query)
	result, err = d.DB.ExecContext(ctx, query)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result.RowsAffected())
	}

	return err
}
