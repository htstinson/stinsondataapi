package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectCustomers(ctx context.Context, schema_name string, limit, offset int) ([]model.Customer, error) {

	fmt.Println("database.go SelectCustomers()")

	query := fmt.Sprintf("SELECT id, name, created_at FROM %s.customers ORDER BY name ASC LIMIT $1 OFFSET $2", schema_name)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing customers: %w", err)
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		if err := rows.Scan(&customer.Id, &customer.Name, &customer.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning customer: %w", err)
		}

		customers = append(customers, customer)
	}
	return customers, nil
}

func (d *Database) CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	fmt.Println("d CreateCustomer")

	profile, err := d.GetProfileByParent(ctx, customer.Subscriber_ID)
	if err != nil {
		fmt.Println("error getting profile")
		fmt.Println(customer.Subscriber_ID)
		return customer, err
	}

	query := fmt.Sprintf(`INSERT INTO %s.customers (id, name, parent_id) VALUES ($1, $2, $3)`, customer.Schema_Name)

	_, err = d.DB.ExecContext(ctx, query,
		customer.Id,
		customer.Name,
		profile.Id,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating customer: %w", err)
	}

	return customer, nil

}
