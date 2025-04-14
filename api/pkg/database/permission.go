package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Permission
func (d *Database) GetPermission(ctx context.Context, id string) (*model.Permission, error) {
	fmt.Println("d GetPermission")

	var permission model.Permission

	err := d.db.QueryRowContext(ctx,
		"SELECT id, name, description, created_at FROM permissions WHERE id = $1",
		id,
	).Scan(&permission.Id, &permission.Name, &permission.Description, &permission.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting permissions: %w", err)
	}

	return &permission, nil
}

func (d *Database) CreatePermission(ctx context.Context, name string, description string) (*model.Permission, error) {
	fmt.Println("d CreatePermission")

	permission := &model.Permission{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}

	query := `
        INSERT INTO permissions (id, name, description, created_at) VALUES ($1, $2, $3, $4)
    `

	_, err := d.db.ExecContext(ctx, query,
		permission.Id,
		permission.Name,
		permission.Description,
		permission.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating permission: %w", err)
	}

	return permission, nil

}

func (d *Database) SelectPermissions(ctx context.Context, limit, offset int) ([]model.Permission, error) {

	fmt.Println("database.go SelectPermissions()")

	rows, err := d.db.QueryContext(ctx,
		"SELECT id, name, description, created_at FROM permissions ORDER BY name ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error selecting permissions: %w", err)
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var permission model.Permission
		if err := rows.Scan(&permission.Id, &permission.Name, &permission.Description, &permission.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning customer: %w", err)
		}

		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (d *Database) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	fmt.Println("d UpdatePermission")

	query := `UPDATE permissinos SET name = $1, description = $2 WHERE id = $3`

	_, err := d.db.ExecContext(ctx, query, permission.Name, permission.Description, permission.Id)

	return err
}

func (d *Database) DeletePermission(ctx context.Context, id string) error {
	fmt.Println("d DeletePermission")

	query := `DELETE FROM permissions WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}
