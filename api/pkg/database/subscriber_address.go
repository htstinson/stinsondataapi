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

	fmt.Println("d SelectSubscriberAddresses")

	if order == "" {
		order = "asc"
	}

	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf(`SELECT id, created_at, modified_at, 
		address_type, address_use, street1, street2, po_box, city, state, zip,
		COUNT(*) OVER() AS total 
		FROM %s.addresses WHERE id = $1 ORDER BY %s %s LIMIT $1 OFFSET $2`, subscriber.Schema_Name, sort, order)

	rows, err := d.DB.QueryContext(ctx,
		query,
		subscriber.Id,
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

	return &addresses, 0, nil
}
