package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectUserSubscriberRoleView(ctx context.Context, user_subscriber_role_view model.User_Subscriber_Role_View, limit, offset int) ([]model.User_Subscriber_Role_View, error) {
	fmt.Println("database.go SelectUserSubscriberRolesView()")

	var rows *sql.Rows
	var query string
	var err error

	// User_ID was not provided
	if user_subscriber_role_view.User_ID == "" {
		query = `SELECT id, user_subscriber_id, role_id, role_name, user_id, username, 
		subscriber_id, subscriber_name, created_at, updated_at 
		FROM common.user_subscriber_role_view 
		ORDER BY username, subscriber_name, role_name ASC LIMIT $1 OFFSET $2`

		rows, err = d.DB.QueryContext(ctx, query, limit, offset)
		if err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error selecting rows: %w", err)
		}
		defer rows.Close()

	} else {
		// User_ID was provided
		query = `SELECT id, user_subscriber_id, role_id, role_name, user_id, username, subscriber_id, subscriber_name, created_at, updated_at 
		FROM common.user_subscriber_role_view 
		ORDER BY username, subscriber_name, role_name where user_id = $1 ASC LIMIT $2 OFFSET $3`

		fmt.Println(query)

		rows, err = d.DB.QueryContext(ctx, query, user_subscriber_role_view.User_ID, limit, offset)
		if err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error selecting rows: %w", err)
		}
		defer rows.Close()
	}

	var user_subscriber_role_views []model.User_Subscriber_Role_View
	for rows.Next() {
		var user_subscriber_role_view model.User_Subscriber_Role_View
		if err := rows.Scan(&user_subscriber_role_view.Id, &user_subscriber_role_view.User_Subscriber_ID,
			&user_subscriber_role_view.Role_Id, &user_subscriber_role_view.Role_Name,
			&user_subscriber_role_view.User_ID, &user_subscriber_role_view.User_Name,
			&user_subscriber_role_view.Subscriber_Id, &user_subscriber_role_view.Subscriber_Name,
			&user_subscriber_role_view.Created_At, &user_subscriber_role_view.Updated_At); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_subscriber_role: %w", err)
		}

		user_subscriber_role_views = append(user_subscriber_role_views, user_subscriber_role_view)
	}
	return user_subscriber_role_views, nil
}

func (d *Database) CreateUserSubscriberRole(ctx context.Context, user_subscriber_id string, role_id string) (*model.User_Subscriber_Role, error) {
	fmt.Println("d CreateUserSubscriberRole")

	user_subscriber_role := &model.User_Subscriber_Role{
		Id:                 uuid.New().String(),
		User_Subscriber_ID: user_subscriber_id,
		Role_Id:            role_id,
	}

	query := `
        INSERT INTO user_subscriber_role (id, user_subscriber_id, role_id) VALUES ($1, $2, $3)
    `

	fmt.Println(query)
	fmt.Println(user_subscriber_role)

	_, err := d.DB.ExecContext(ctx, query,
		user_subscriber_role.Id,
		user_subscriber_role.User_Subscriber_ID,
		user_subscriber_role.Role_Id,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user_subscriber_role: %w", err)
	}

	return user_subscriber_role, nil

}

func (d *Database) LookupUserSubscriberRole(ctx context.Context, user_subscriber_id string, role_id string) (*model.User_Subscriber_Role, error) {
	fmt.Println("d LookupUserSubscriberRole")

	var user_subscriber_role = model.User_Subscriber_Role{}

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, user_subscriber_id, role_id FROM user_subscriber_role WHERE user_subscriber_id = $1 and role_id = $2",
		user_subscriber_id, role_id,
	).Scan(&user_subscriber_role.Id, &user_subscriber_role.User_Subscriber_ID, &user_subscriber_role.Role_Id)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user_subscriber_role: %w", err)
	}

	return &user_subscriber_role, nil
}

func (d *Database) UpdateUserSubscriberRole(ctx context.Context, user_subscriber_role model.User_Subscriber_Role) error {
	fmt.Println("d UpdateUserSubscriberRole")

	query := `UPDATE user_subscriber_role SET user_subscriber_id = $1, role_id = $2 WHERE id = $3`
	fmt.Println("id", user_subscriber_role.Id)
	fmt.Println("user_subscriber_id", user_subscriber_role.User_Subscriber_ID)
	fmt.Println("role id", user_subscriber_role.Role_Id)

	_, err := d.DB.ExecContext(ctx, query, user_subscriber_role.User_Subscriber_ID, user_subscriber_role.Role_Id, user_subscriber_role.Id)

	return err
}

func (d *Database) GetUserSubscriberRole(ctx context.Context, id string) (*model.User_Subscriber_Role, error) {
	fmt.Println("d GetUserSubscriberRole")

	var user_subscriber_role model.User_Subscriber_Role

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, user_subscriber_id, role_id FROM user_subscriber_role WHERE id = $1",
		id,
	).Scan(&user_subscriber_role.Id, &user_subscriber_role.User_Subscriber_ID, &user_subscriber_role.Role_Id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting subscriber: %w", err)
	}

	return &user_subscriber_role, nil
}

func (d *Database) DeleteUserSubscriberRole(ctx context.Context, id string) error {
	fmt.Println("d DeleteUserSubscriberRole")

	query := `DELETE FROM user_subscriber_role WHERE id = $1`

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}
