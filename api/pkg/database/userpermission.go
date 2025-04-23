package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectUserPermissions(ctx context.Context, limit, offset int) ([]model.User_Permission, error) {
	fmt.Println("database.go SelectUserPermission()")
	rows, err := d.DB.QueryContext(ctx,
		"SELECT user_id, permission_id, created_at FROM user_permissions ORDER BY user_id ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var user_permissions []model.User_Permission
	for rows.Next() {
		var user_permission model.User_Permission
		if err := rows.Scan(&user_permission.User_Id, &user_permission.Permission_Id, &user_permission.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_permission: %w", err)
		}

		user_permissions = append(user_permissions, user_permission)
	}
	return user_permissions, nil
}
