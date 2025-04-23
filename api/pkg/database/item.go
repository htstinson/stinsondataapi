package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Item
func (d *Database) GetItem(ctx context.Context, id string) (*model.Item, error) {
	var item model.Item
	err := d.DB.QueryRowContext(ctx,
		"SELECT id, name, created_at FROM items WHERE id = $1",
		id,
	).Scan(&item.ID, &item.Name, &item.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting item: %w", err)
	}
	return &item, nil
}

func (d *Database) UpdateItem(ctx context.Context, item *model.Item) error {

	query := `UPDATE items SET name = $1 WHERE id = $2`

	_, err := d.DB.ExecContext(ctx, query, item.Name, item.ID)

	return err

}

func (d *Database) DeleteItem(ctx context.Context, id string) error {

	query := `DELETE FROM items WHERE id = $1`

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}

func (d *Database) CreateItem(ctx context.Context, item *model.Item) error {
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()

	_, err := d.DB.ExecContext(ctx,
		"INSERT INTO items (id, name, created_at) VALUES ($1, $2, $3)",
		item.ID, item.Name, item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating item: %w", err)
	}
	return nil
}

func (d *Database) SelectItems(ctx context.Context, limit, offset int) ([]model.Item, error) {
	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, name, created_at FROM items ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
