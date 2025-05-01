package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

//Roles

func (d *Database) SelectRolesByUser(ctx context.Context, userId string) (model.Roles, error) {
	fmt.Println("d SelectRolesByUser")

	var roles = model.Roles{}

	query := `
        SELECT user_id, username, role_name
        FROM user_customer_roles_view
        WHERE user_id = $1
    `

	err := d.DB.QueryRowContext(ctx, query, userId).Scan(
		&roles.Id, &roles.Username, &roles.Names,
	)

	return roles, err
}

func (d *Database) SelectRoles(ctx context.Context, limit, offset int) ([]model.Role, error) {
	fmt.Println("d SelectRoles")

	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, name FROM roles ORDER BY name ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.Id, &role.Name); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user: %w", err)
		}

		roles = append(roles, role)
	}

	return roles, err
}

func (d *Database) GetRole(ctx context.Context, id string) (*model.Role, error) {
	fmt.Println("d GetRole")

	var role model.Role

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, name, created_at FROM roles WHERE id = $1",
		id,
	).Scan(&role.Id, &role.Name, &role.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	return &role, nil
}

func (d *Database) UpdateRole(ctx context.Context, role *model.Role) error {
	fmt.Println("d UpdateRole")

	query := `UPDATE roles SET name = $1 WHERE id = $2`

	_, err := d.DB.ExecContext(ctx, query, role.Name, role.Id)

	return err
}

func (d *Database) CreateRole(ctx context.Context, name string) (*model.Role, error) {
	fmt.Println("d CreateRole")

	role := &model.Role{
		Id:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO roles (id, name, created_at) VALUES ($1, $2, $3)
    `

	_, err := d.DB.ExecContext(ctx, query, role.Id, role.Name, role.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating role: %w", err)
	}

	return role, nil
}

func (d *Database) DeleteRole(ctx context.Context, id string) error {
	fmt.Println("d DeleteRole")

	query := `DELETE FROM roles WHERE id = $1`

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}
