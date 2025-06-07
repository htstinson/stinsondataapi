package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectCustomers(ctx context.Context, subscriber model.Subscriber, limit, offset int) ([]model.Customer, error) {

	fmt.Println("database.go SelectCustomers()")

	query := fmt.Sprintf("SELECT id, name, created_at FROM %s.customers ORDER BY name ASC LIMIT $1 OFFSET $2", subscriber.Schema_Name)

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

		customer.Schema_Name = subscriber.Schema_Name
		customer.Subscriber_ID = subscriber.Id

		customers = append(customers, customer)
	}
	return customers, nil
}

func (d *Database) CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	fmt.Println("d CreateCustomer")

	subcriber, err := d.GetSubscriber(ctx, customer.Subscriber_ID)
	if err != nil {
		fmt.Println(err.Error())
		return customer, err
	}

	profile, err := d.GetProfileByParent(ctx, subcriber)
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

func (d *Database) GetCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	fmt.Println("d GetCustomer")

	query := fmt.Sprintf(`SELECT name, subscriber_id, schema_name, created_at FROM %s.customer WHERE id = $1`, customer.Schema_Name)

	fmt.Println(query)

	err := d.DB.QueryRowContext(ctx, query, customer.Id).Scan(&customer.Name, &customer.Subscriber_ID, &customer.Schema_Name, &customer.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}
	return customer, nil
}

func (d *Database) DeleteCustomer(ctx context.Context, customer *model.Customer) error {

	fmt.Println("d DeleteCustomer")

	query := fmt.Sprintf(`DELETE FROM %s.customers WHERE id = $1`, customer.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query, customer.Id)

	return err
}
