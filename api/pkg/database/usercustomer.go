package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectUserCustomerView(ctx context.Context, limit, offset int) ([]model.User_Customer_View, error) {
	fmt.Println("database.go SelectUserCustomerView()")
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, user_id, customer_id, user_username, customer_name FROM user_customer_view ORDER BY user_username, customer_name ASC LIMIT $1 OFFSET $2",
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
		if err := rows.Scan(&user_customer_view.Id, &user_customer_view.User_ID, &user_customer_view.Customer_Id, &user_customer_view.User_Username, &user_customer_view.Customer_Name); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_permission: %w", err)
		}

		user_customer_views = append(user_customer_views, user_customer_view)
	}
	return user_customer_views, nil
}

func (d *Database) UpdateUserCustomer(ctx context.Context, user_customer model.User_Customer) error {
	fmt.Println("d UpdateUserCustomer")

	query := `UPDATE user_customer SET user_id = $1, customer_id = $2 WHERE id = $3`
	fmt.Println(user_customer.User_ID)
	fmt.Println(user_customer.Customer_Id)
	fmt.Println(user_customer.Id)

	_, err := d.db.ExecContext(ctx, query, user_customer.User_ID, user_customer.Customer_Id, user_customer.Id)

	return err
}

func (d *Database) GetUserCustomer(ctx context.Context, id string) (*model.User_Customer, error) {
	fmt.Println("d GetUserCustomer")

	var user_customer model.User_Customer

	err := d.db.QueryRowContext(ctx,
		"SELECT id, user_id, customer_id FROM user_customer WHERE id = $1",
		id,
	).Scan(&user_customer.Id, &user_customer.Id, &user_customer.Id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	return &user_customer, nil
}
