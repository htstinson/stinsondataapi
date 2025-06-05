package database

import (
	"context"
	"fmt"

	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectContacts(ctx context.Context, customer model.Customer, limit, offset int) ([]model.Contact, error) {

	fmt.Println("database.go SelectContacts()")

	query := fmt.Sprintf(`SELECT id, parent_id, lastname, firstname, created_at FROM %s.contacts 
	WHERE parent_id = '%s' ORDER BY lastname,firstname ASC LIMIT $1 OFFSET $2`, customer.Schema_Name, customer.Id)

	rows, err := d.DB.QueryContext(ctx,
		query,
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing contacts: %w", err)
	}
	defer rows.Close()

	var contacts []model.Contact
	for rows.Next() {
		var contact model.Contact
		if err := rows.Scan(&contact.Id, &contact.ParentId, &contact.LastName, &contact.FirstName, &contact.CreatedAt); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning contact: %w", err)
		}

		contacts = append(contacts, contact)
	}
	return contacts, nil
}
