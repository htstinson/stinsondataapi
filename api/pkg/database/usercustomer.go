package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectUserCustomerView(ctx context.Context, limit, offset int) ([]model.User_Customer_View, error) {
	fmt.Println("database.go SelectUserCustomerView()")
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, user_id, customer_id, user_username, customer_name, assigned_at FROM user_customer_view ORDER BY user_username, customer_name ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var user_customer_views []model.User_Customer_View
	for rows.Next() {
		var user_customer_view model.User_Customer_View
		if err := rows.Scan(&user_customer_view.Id, &user_customer_view.User_ID, &user_customer_view.Customer_Id, &user_customer_view.User_Username, &user_customer_view.Customer_Name, &user_customer_view.Assignedd_At); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_permission: %w", err)
		}

		user_customer_views = append(user_customer_views, user_customer_view)
	}
	return user_customer_views, nil
}
