package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectRolePermissions(ctx context.Context, limit, offset int) ([]model.Role_Permission, error) {
	fmt.Println("database.go SelectRolePermission()")
	rows, err := d.db.QueryContext(ctx,
		"SELECT role_id, permission_id, created_at FROM Role_permissions ORDER BY Role_id ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var role_permissions []model.Role_Permission
	for rows.Next() {
		var role_permission model.Role_Permission
		if err := rows.Scan(&role_permission.Role_Id, &role_permission.Permission_Id, &role_permission.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_permission: %w", err)
		}

		role_permissions = append(role_permissions, role_permission)
	}
	return role_permissions, nil
}
