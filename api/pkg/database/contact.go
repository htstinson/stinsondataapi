package database

import (
	"context"
	"database/sql"
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
		if err == sql.ErrNoRows {
			return nil, err
		}
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

		contact.Schema_Name_ = customer.Schema_Name

		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (d *Database) CreateContact(ctx context.Context, contact *model.Contact) (*model.Contact, error) {
	fmt.Println("d CreateCustomer")

	query := fmt.Sprintf(`INSERT INTO %s.contacts (parent_id, lastname, firstname) VALUES ($1, $2, $3)`, contact.Schema_Name_)

	_, err := d.DB.ExecContext(ctx, query, contact.ParentId, contact.LastName, contact.FirstName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating customer: %w", err)
	}

	return contact, nil

}

func (d *Database) DeleteContact(ctx context.Context, contact *model.Contact) error {

	query := fmt.Sprintf(`DELETE FROM %s.contacts WHERE id = $1`, contact.Schema_Name_)

	_, err := d.DB.ExecContext(ctx, query, contact.Id)

	return err
}

func (d *Database) GetContact(ctx context.Context, c model.Contact) (*model.Contact, error) {
	fmt.Println("d GetContact")

	query := fmt.Sprintf(`SELECT parent_id, lastname, firstname, created_at FROM %s.contacts WHERE id = $1`, c.Schema_Name_)

	fmt.Println(query)

	contact := &model.Contact{
		Id: c.Id,
	}

	err := d.DB.QueryRowContext(ctx, query, c.Id).Scan(&contact.ParentId, &contact.LastName, &contact.FirstName, &contact.CreatedAt)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting contact: %w", err)
	}
	return contact, nil
}

func (d *Database) UpdateContact(ctx context.Context, contact *model.Contact) error {
	fmt.Println("d UpdateContact")

	query := fmt.Sprintf(`UPDATE %s.contacts SET lastname = $2, firstname = $3 WHERE id = $1`, contact.Schema_Name_)
	fmt.Println(query)
	fmt.Println(contact.Id)
	fmt.Println(contact.LastName)
	fmt.Println(contact.FirstName)

	_, err := d.DB.ExecContext(ctx, query, contact.Id, contact.LastName, contact.FirstName)
	if err != nil {
		fmt.Println(err.Error())
	}

	return err

}
