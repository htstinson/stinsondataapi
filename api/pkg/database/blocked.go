package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Admin - Blocked
func (d *Database) SelectBlocked(ctx context.Context, limit int, offset int, sort string, order string) ([]model.Blocked, error) {

	q := fmt.Sprintf("SELECT id, ip, notes, created_at FROM blocked ORDER BY %s %s LIMIT $1 OFFSET $2", sort, strings.ToUpper(order))

	rows, err := d.db.QueryContext(ctx, q, limit, offset)

	if err != nil {
		fmt.Printf("[%v] [database][SelectBlocked] error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
		fmt.Printf("[%v] [database][SelectBlocked] query: %s.\n", time.Now().Format(time.RFC3339), q)
		fmt.Printf("[%v] [database][SelectBlocked] limit %v offset %v sort %s order %s.\n", time.Now().Format(time.RFC3339), limit, offset, strings.ToUpper(sort), strings.ToUpper(order))
		return nil, fmt.Errorf("query error")
	}
	defer rows.Close()

	var items []model.Blocked

	for rows.Next() {
		var item model.Blocked
		var notesNullable sql.NullString
		if err := rows.Scan(&item.ID, &item.IP, &notesNullable, &item.CreatedAt); err != nil {
			fmt.Printf("[%v] [database][ListBlocked] error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
			return nil, fmt.Errorf("error scanning blocked: %w", err)
		}
		if notesNullable.Valid {
			item.Notes = notesNullable.String
		} else {
			item.Notes = ""
		}

		items = append(items, item)
	}

	return items, nil
}

func (d *Database) UpdateBlocked(ctx context.Context, blocked *model.Blocked) error {
	fmt.Println("d UpdateBlocked", blocked.IP, blocked.Notes)

	query := `UPDATE blocked SET ip=$1, notes=$2 WHERE id = $3`

	_, err := d.db.ExecContext(ctx, query, blocked.IP, blocked.Notes, blocked.ID)

	return err

}

func (d *Database) GetBlockedByIP(ctx context.Context, ip string) (*model.Blocked, error) {
	var blocked model.Blocked
	var notesNull sql.NullString // Use sql.NullString to handle NULL values

	query := fmt.Sprintf(`SELECT id, ip, notes, created_at FROM blocked WHERE ip ='%s'`, ip)

	err := d.db.QueryRowContext(ctx, query).Scan(&blocked.ID, &blocked.IP, &notesNull, &blocked.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		fmt.Println("d GetBlockedByIP", err.Error())
		return nil, fmt.Errorf("error getting blocked: %w", err)
	}

	// Only assign the value if it's not NULL
	if notesNull.Valid {
		blocked.Notes = notesNull.String
	} else {
		blocked.Notes = "" // Or another default value of your choice
	}

	return &blocked, nil
}

func (d *Database) GetBlocked(ctx context.Context, id string) (*model.Blocked, error) {
	var blocked model.Blocked
	var notesNull sql.NullString // Use sql.NullString to handle NULL values

	query := "SELECT id, ip, notes, created_at FROM blocked WHERE id = $1"

	err := d.db.QueryRowContext(ctx, query, id).Scan(&blocked.ID, &blocked.IP, &notesNull, &blocked.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		fmt.Println("d GetBlocked", err.Error())
		return nil, fmt.Errorf("error getting blocked: %w", err)
	}

	// Only assign the value if it's not NULL
	if notesNull.Valid {
		blocked.Notes = notesNull.String
	} else {
		blocked.Notes = "" // Or another default value of your choice
	}

	return &blocked, nil
}

func (d *Database) CreateBlocked(ctx context.Context, blocked model.Blocked) (*model.Blocked, error) {
	fmt.Println("d CreateBlocked(ctx, blocked)")
	blocked.CreatedAt = time.Now()

	_, err := d.GetBlockedByIP(ctx, blocked.IP)
	if err == nil {
		return nil, errors.New("duplicate")
	}

	query := `
        INSERT INTO blocked (ip, notes, created_at)
        VALUES ($1, $2, $3)
    `

	_, err = d.db.ExecContext(ctx, query, blocked.IP, blocked.Notes, blocked.CreatedAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating blocked: %w", err)
	}

	return &blocked, nil

}

func (d *Database) DeleteBlocked(ctx context.Context, id string) error {

	query := `DELETE FROM blocked WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}
