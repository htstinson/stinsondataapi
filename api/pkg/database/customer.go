package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectCustomers(ctx context.Context, subscriber model.Subscriber, limit int, offset int, sort string, order string) ([]model.Customer, int, error) {

	fmt.Println("d SelectCustomers")

	if order == "" {
		order = "asc"
	}

	if sort == "" {
		sort = "name"
	}

	query := fmt.Sprintf("SELECT id, name, created_at, COUNT(*) OVER() AS total FROM %s.customers ORDER BY %s %s LIMIT $1 OFFSET $2", subscriber.Schema_Name, sort, order)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, 0, fmt.Errorf("error listing customers: %w", err)
	}
	defer rows.Close()

	var total int
	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		if err := rows.Scan(&customer.Id, &customer.Name, &customer.CreatedAt, &total); err != nil {
			fmt.Println(err.Error())
			return nil, 0, fmt.Errorf("error scanning customer: %w", err)
		}

		customer.Schema_Name = subscriber.Schema_Name
		customer.Subscriber_Id = subscriber.Id

		customers = append(customers, customer)
	}
	return customers, total, nil
}

func (d *Database) CreateCustomer(ctx context.Context, customer *model.Customer, profile *model.Profile) (*model.Customer, error) {
	fmt.Println("d CreateCustomer")

	subcriber, err := d.GetSubscriber(ctx, customer.Subscriber_Id)
	if err != nil {
		fmt.Println(err.Error())
		return customer, err
	}

	subcriber.Id = customer.Subscriber_Id

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

func (d *Database) GetCustomer(ctx context.Context, temp_customer model.Customer) (*model.Customer, error) {
	fmt.Println("d GetCustomer")

	query := fmt.Sprintf(`SELECT name, created_at FROM %s.customers WHERE id = $1`, temp_customer.Schema_Name)

	customer := &model.Customer{
		Id:            temp_customer.Id,
		Subscriber_Id: temp_customer.Subscriber_Id,
		Schema_Name:   temp_customer.Schema_Name,
	}

	err := d.DB.QueryRowContext(ctx, query, temp_customer.Id).Scan(&customer.Name, &customer.CreatedAt)

	if err == sql.ErrNoRows {
		fmt.Println(err.Error())
		return nil, nil
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	return customer, nil
}

func (d *Database) DeleteCustomer(ctx context.Context, customer *model.Customer) error {

	fmt.Println("d DeleteCustomer")

	query := fmt.Sprintf(`DELETE FROM %s.customers WHERE id = $1`, customer.Schema_Name)
	fmt.Println(query)
	fmt.Println(customer)

	_, err := d.DB.ExecContext(ctx, query, customer.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

func (d *Database) UpdateCustomer(ctx context.Context, customer *model.Customer) error {
	fmt.Println("d UpdateCustomer")

	query := fmt.Sprintf(`UPDATE %s.customers SET name = $2 WHERE id = $1`, customer.Schema_Name)

	fmt.Println(query)
	fmt.Println(customer.Name)

	_, err := d.DB.ExecContext(ctx, query, customer.Id, customer.Name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
