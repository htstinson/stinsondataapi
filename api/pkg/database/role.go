package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

//Roles

func (d *Database) SelectRolesByUser(ctx context.Context, userId string) (model.Roles, error) {
	fmt.Println("d SelectRoles")

	var roles = model.Roles{}

	query := `
        SELECT user_id, username, role_name
        FROM user_roles_view
        WHERE user_id = $1
    `

	err := d.db.QueryRowContext(ctx, query, userId).Scan(
		&roles.Id, &roles.Username, &roles.Names,
	)

	return roles, err
}

func (d *Database) SelectRoles(ctx context.Context, limit, offset int) ([]model.Role, error) {
	fmt.Println("d SelectRoles")

	rows, err := d.db.QueryContext(ctx,
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
