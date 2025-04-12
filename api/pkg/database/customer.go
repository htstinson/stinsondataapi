package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Customers
func (d *Database) GetCustomer(ctx context.Context, id string) (*model.Customer, error) {
	fmt.Println("d GetCustomer")

	var customer model.Customer

	err := d.db.QueryRowContext(ctx,
		"SELECT id, name, created_at FROM customers WHERE id = $1",
		id,
	).Scan(&customer.ID, &customer.Name, &customer.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	return &customer, nil
}

func (d *Database) GetCustomerByName(ctx context.Context, name string) (*model.Customer, error) {
	fmt.Printf("[%v] [GetCustomerByName] %s\n", time.Now().Format(time.RFC3339), name)

	customer := &model.Customer{}
	query := `
        SELECT id, name, created_at FROM customers WHERE username = $1
    `

	err := d.db.QueryRowContext(ctx, query, name).Scan(
		&customer.ID,
		&customer.Name,
		&customer.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	return customer, nil
}

func (d *Database) CreateCustomer(ctx context.Context, name string) (*model.Customer, error) {
	fmt.Println("CreateCustomer")

	customer := &model.Customer{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO customers (id, name, created_at) VALUES ($1, $2, $3)
    `

	_, err := d.db.ExecContext(ctx, query,
		customer.ID,
		customer.Name,
		customer.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating customer: %w", err)
	}

	return customer, nil

}

func (d *Database) SelectCustomers(ctx context.Context, limit, offset int) ([]model.Customer, error) {

	fmt.Println("database.go SelectCustomers()")

	rows, err := d.db.QueryContext(ctx,
		"SELECT id, name, created_at FROM customers ORDER BY name ASC LIMIT $1 OFFSET $2",
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
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning customer: %w", err)
		}

		customers = append(customers, customer)
	}
	return customers, nil
}

func (d *Database) UpdateCustomer(ctx context.Context, customer *model.Customer) error {
	fmt.Println("h UpdateCustomer")

	query := `UPDATE customers SET name = $1 WHERE id = $2`

	_, err := d.db.ExecContext(ctx, query, customer.Name, customer.ID)

	return err
}

func (d *Database) DeleteCustomer(ctx context.Context, id string) error {
	query := `DELETE FROM customers WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}
