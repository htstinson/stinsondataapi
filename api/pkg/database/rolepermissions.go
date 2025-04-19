package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectRolePermissionsView(ctx context.Context, limit, offset int) ([]model.Role_Permission_View, error) {
	fmt.Println("database.go SelectRolePermissionsView()")

	rows, err := d.db.QueryContext(ctx,
		"SELECT role_id, role_name, permission_id, permission_name, object_id, object_name, object_type, created_at FROM role_permissions_view ORDER BY Role_name, permission_name, object_name ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var role_permissions_view []model.Role_Permission_View
	for rows.Next() {
		var role_permission_view model.Role_Permission_View
		if err := rows.Scan(&role_permission_view.Role_Id, &role_permission_view.V_Role_Name,
			&role_permission_view.Permission_Id, &role_permission_view.V_Permission_Name,
			&role_permission_view.Object_Id, &role_permission_view.V_Object_Name, &role_permission_view.V_Object_Type,
			&role_permission_view.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_permission: %w", err)
		}

		role_permissions_view = append(role_permissions_view, role_permission_view)
	}
	return role_permissions_view, nil
}
