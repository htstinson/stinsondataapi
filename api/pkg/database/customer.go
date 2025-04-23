package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Customers
func (d *Database) GetCustomer(ctx context.Context, id string) (*model.Customer, error) {
	fmt.Println("d GetCustomer")

	var customer model.Customer

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, name, created_at FROM customers WHERE id = $1",
		id,
	).Scan(&customer.Id, &customer.Name, &customer.CreatedAt)

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

	err := d.DB.QueryRowContext(ctx, query, name).Scan(
		&customer.Id,
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
	fmt.Println("d CreateCustomer")

	customer := &model.Customer{
		Id:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO customers (id, name, created_at) VALUES ($1, $2, $3)
    `

	_, err := d.DB.ExecContext(ctx, query,
		customer.Id,
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

	rows, err := d.DB.QueryContext(ctx,
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
		if err := rows.Scan(&customer.Id, &customer.Name, &customer.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning customer: %w", err)
		}

		customers = append(customers, customer)
	}
	return customers, nil
}

func (d *Database) UpdateCustomer(ctx context.Context, customer *model.Customer) error {
	fmt.Println("d UpdateCustomer")

	query := `UPDATE customers SET name = $1 WHERE id = $2`

	_, err := d.DB.ExecContext(ctx, query, customer.Name, customer.Id)

	return err
}

func (d *Database) DeleteCustomer(ctx context.Context, id string) error {
	fmt.Println("d DeleteCustomer")

	query := `DELETE FROM customers WHERE id = $1`

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}

func (d *Database) Create_Schema(ctx context.Context, new_schema_name string) error {
	fmt.Println("d Create_Schema")

	err := errors.New("temporary stop")

	return err

	// err := Copy_Schema(d.db, d.Config, false, new_schema_name)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// return err
}
