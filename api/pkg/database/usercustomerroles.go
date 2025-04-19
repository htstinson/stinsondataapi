package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectUserCustomerRolesView(ctx context.Context, limit, offset int) ([]model.User_Customer_Roles_View, error) {
	fmt.Println("database.go SelectUserCustomerRolesView()")
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, user_customer_id, role_id, role_name, user_id, user_username, customer_id, customer_name, created_at, updated_at, FROM user_customer_roles_view ORDER BY user_username, customer_name, role_name ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var user_customer_roles_views []model.User_Customer_Roles_View
	for rows.Next() {
		var user_customer_roles_view model.User_Customer_Roles_View
		if err := rows.Scan(&user_customer_roles_view.Id,
			&user_customer_roles_view.Role_Id, &user_customer_roles_view.Role_Name,
			&user_customer_roles_view.User_ID, &user_customer_roles_view.User_Name,
			&user_customer_roles_view.Customer_Id, &user_customer_roles_view.Customer_Name,
			&user_customer_roles_view.Created_At, &user_customer_roles_view.Updated_At); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_customer_role: %w", err)
		}

		user_customer_roles_views = append(user_customer_roles_views, user_customer_roles_view)
	}
	return user_customer_roles_views, nil
}
