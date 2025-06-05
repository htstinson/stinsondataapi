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

func (d *Database) CreateContact(ctx context.Context, contact *model.Contact) (*model.Contact, error) {
	fmt.Println("d CreateCustomer")

	query := fmt.Sprintf(`INSERT INTO %s.contacts (parent_id, last_name, first_name) VALUES ($1, $2, $3)`, contact.Schema_Name_)

	_, err := d.DB.ExecContext(ctx, query, contact.ParentId, contact.LastName, contact.FirstName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating customer: %w", err)
	}

	return contact, nil

}
