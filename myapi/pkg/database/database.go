package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"myapi/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository interface {
	GetItem(ctx context.Context, id string) (*model.Item, error)
	CreateItem(ctx context.Context, item *model.Item) error
	ListItems(ctx context.Context, limit, offset int) ([]model.Item, error)
	Close() error
}

type Database struct {
	db *sql.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func New(cfg Config) (Repository, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Initialize schema
	if err := initializeSchema(db); err != nil {
		return nil, fmt.Errorf("error initializing schema: %w", err)
	}

	return &Database{db: db}, nil
}

func initializeSchema(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            id VARCHAR(36) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL
        );
        CREATE INDEX IF NOT EXISTS items_created_at_idx ON items(created_at DESC);
    `)
	return err
}

func (d *Database) GetItem(ctx context.Context, id string) (*model.Item, error) {
	var item model.Item
	err := d.db.QueryRowContext(ctx,
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

func (d *Database) CreateItem(ctx context.Context, item *model.Item) error {
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()

	_, err := d.db.ExecContext(ctx,
		"INSERT INTO items (id, name, created_at) VALUES ($1, $2, $3)",
		item.ID, item.Name, item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating item: %w", err)
	}
	return nil
}

func (d *Database) ListItems(ctx context.Context, limit, offset int) ([]model.Item, error) {
	rows, err := d.db.QueryContext(ctx,
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

func (d *Database) Close() error {
	return d.db.Close()
}
