package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Background
func (d *Database) SelectSubscriberBackgrounds(ctx context.Context, subscriber model.Subscriber,
	limit int, offset int, sort string, order string) (*[]model.Background, int, error) {

	if order == "" {
		order = "asc"
	}

	if sort == "" {
		sort = "topic"
	}

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, 
		topic, summary, details,
		COUNT(*) OVER() AS total 
		FROM %s.background ORDER BY %s %s LIMIT $1 OFFSET $2`, subscriber.Schema_Name, sort, order)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, 0, fmt.Errorf("error listing backgrounds: %w", err)
	}
	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	var total int
	var backgrounds []model.Background
	for rows.Next() {
		var background model.Background
		if err := rows.Scan(&background.Id, &background.CreatedAt, &background.ModifiedAt,
			&background.Topic, &background.Summary, &background.Details,
			&total); err != nil {
			fmt.Println(err.Error())
			return nil, 0, fmt.Errorf("error scanning background: %w", err)
		}

		background.SubscriberId = subscriber.Id

		backgrounds = append(backgrounds, background)
	}

	return &backgrounds, total, nil
}

func (d *Database) GetSubscriberBackground(ctx context.Context, subscriber_schema_name string, background_id string) (*model.Background, error) {
	fmt.Println("d Get Subscriber Background")

	query := fmt.Sprintf(`SELECT id, subscriber_id, created_at, modified_at, 
		topic, summary, details FROM %s.background WHERE id = $1`, subscriber_schema_name)

	rows, err := d.DB.QueryContext(ctx, query, background_id)
	if err != nil {
		return nil, fmt.Errorf("error listing backgrounds: %w", err)
	}
	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil, nil
	}

	var background = model.Background{}
	for rows.Next() {
		if err := rows.Scan(&background.Id, &background.SubscriberId, &background.CreatedAt, &background.ModifiedAt,
			&background.Topic, &background.Summary, &background.Details); err != nil {
			return nil, fmt.Errorf("error scanning background: %w", err)
		}
	}
	return &background, nil
}

func (d *Database) UpdateSubscriberBackground(ctx context.Context, subscriber *model.Subscriber, background model.Background) error {
	fmt.Println("d UpdateSubscriberBackground")

	query := fmt.Sprintf(`UPDATE %s.background SET 
	topic = $1, 
	summary = $2, 
	details = $3 
	WHERE id = $4`, subscriber.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query,
		background.Topic,
		background.Summary,
		background.Details,
		background.Id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (d *Database) CreateSubscriberBackground(ctx context.Context, subscriber *model.Subscriber, background model.Background) error {
	fmt.Println("d Create Subscriber Background")

	query := fmt.Sprintf(`INSERT INTO %s.background (subscriber_id, topic, summary, details) 
	values ($1, $2, $3, $4)`, subscriber.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query,
		subscriber.Id,
		background.Topic,
		background.Summary,
		background.Details)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (d *Database) DeleteSubscriberBackground(ctx context.Context, subscriber_schema_name string, background_id string) error {
	fmt.Println("d DeleteSubscriberBackground")

	query := fmt.Sprintf(`DELETE FROM %s.background WHERE id = $1`, subscriber_schema_name)

	_, err := d.DB.ExecContext(ctx, query, background_id)

	return err
}
