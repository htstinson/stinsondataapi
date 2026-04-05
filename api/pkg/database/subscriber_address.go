package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

// Addresses
func (d *Database) SelectSubscriberAddresses(ctx context.Context, subscriber model.Subscriber,
	limit int, offset int, sort string, order string) (*[]model.Address, int, error) {

	if order == "" {
		order = "asc"
	}

	if sort == "" {
		sort = "address_use"
	}

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, 
		address_type, address_use, street1, street2, po_box, city, state, zip,
		COUNT(*) OVER() AS total 
		FROM %s.addresses ORDER BY %s %s LIMIT $1 OFFSET $2`, subscriber.Schema_Name, sort, order)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, 0, fmt.Errorf("error listing addresses: %w", err)
	}
	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil, 0, nil
	}

	var total int
	var addresses []model.Address
	for rows.Next() {
		var address model.Address
		if err := rows.Scan(&address.Id, &address.CreatedAt, &address.ModifiedAt,
			&address.AddressType, &address.AddressUse, &address.Street1, &address.Street2, &address.POBox, &address.City, &address.State, &address.Zip,
			&total); err != nil {
			fmt.Println(err.Error())
			return nil, 0, fmt.Errorf("error scanning address: %w", err)
		}

		address.SubscriberId = subscriber.Id

		addresses = append(addresses, address)
	}

	return &addresses, total, nil
}

func (d *Database) GetSubscriberAddress(ctx context.Context, subscriber_schema_name string, address_id string) (*model.Address, error) {
	fmt.Println("d Get Subscriber Address")

	query := fmt.Sprintf(`SELECT id, subscriber_id, created_at, modified_at, 
		address_type, address_use, street1, street2, po_box, 
		city, state, zip FROM %s.addresses WHERE id = $1`, subscriber_schema_name)

	rows, err := d.DB.QueryContext(ctx, query, address_id)
	if err != nil {
		return nil, fmt.Errorf("error listing addresses: %w", err)
	}
	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil, nil
	}

	var address = model.Address{}
	for rows.Next() {
		if err := rows.Scan(&address.Id, &address.SubscriberId, &address.CreatedAt, &address.ModifiedAt,
			&address.AddressType, &address.AddressUse, &address.Street1, &address.Street2, &address.POBox,
			&address.City, &address.State, &address.Zip); err != nil {

			return nil, fmt.Errorf("error scanning address: %w", err)
		}
	}
	return &address, nil
}

func (d *Database) UpdateSubscriberAddress(ctx context.Context, subscriber *model.Subscriber, address model.Address) error {
	fmt.Println("d UpdateSubscriberAddress")

	query := fmt.Sprintf(`UPDATE %s.addresses SET 
	address_type = $1, 
	address_use = $2, 
	street1 = $3, 
	street2 = $4, 
	po_box = $5, 
	city = $6, 
	state = $7, 
	zip = $8 
	WHERE id = $9`, subscriber.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query,
		address.AddressType,
		address.AddressUse,
		address.Street1,
		address.Street2,
		address.POBox,
		address.City,
		address.State,
		address.Zip,
		address.Id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (d *Database) CreateSubscriberAddress(ctx context.Context, subscriber *model.Subscriber, address model.Address) error {
	fmt.Println("d Create Subscriber Address")

	query := fmt.Sprintf(`INSERT INTO %s.addresses (subscriber_id, address_type, 
	address_use, street1, street2, po_box, city, state, zip) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, subscriber.Schema_Name)

	_, err := d.DB.ExecContext(ctx, query,
		subscriber.Id,
		address.AddressType,
		address.AddressUse,
		address.Street1,
		address.Street2,
		address.POBox,
		address.City,
		address.State,
		address.Zip)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (d *Database) DeleteSubscriberAddress(ctx context.Context, subscriber_schema_name string, address_id string) error {
	fmt.Println("d DeleteSubscriberAddress")

	query := fmt.Sprintf(`DELETE FROM %s.addresses WHERE id = $1`, subscriber_schema_name)

	_, err := d.DB.ExecContext(ctx, query, address_id)

	return err
}
