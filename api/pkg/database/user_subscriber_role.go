package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectUserSubscriberRolesView(ctx context.Context, limit, offset int) ([]model.User_Subscriber_Roles_View, error) {
	fmt.Println("database.go SelectUserSubscriberRolesView()")
	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, user_subscriber_id, role_id, role_name, user_id, username, subscriber_id, subscriber_name, created_at, updated_at FROM user_subscriber_role_view ORDER BY username, subscriber_name, role_name ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var user_subscriber_role_views []model.User_Subscriber_Roles_View
	for rows.Next() {
		var user_subscriber_role_view model.User_Subscriber_Roles_View
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
	fmt.Println("d CreateUserSubscriber")

	user_subscriber_role := &model.User_Subscriber_Role{
		Id:                 uuid.New().String(),
		User_Subscriber_ID: user_subscriber_id,
		Role_Id:            role_id,
	}

	query := `
        INSERT INTO user_subscriber_role (id, user_subscriber_id, role_id) VALUES ($1, $2, $3)
    `

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
