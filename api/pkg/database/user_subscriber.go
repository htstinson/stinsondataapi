package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

func (d *Database) SelectUserSubscriberView(ctx context.Context, user_id string, limit int, offset int) ([]model.User_Subscriber_View, error) {
	fmt.Println("database.go SelectUserSubscriberView()")

	where_clause := ""

	if user_id != "" {
		_, err := ValidateUUID(user_id)
		if err != nil {
			fmt.Printf("Invalid UUID error: %v\n", err)
			return nil, err
		} else {
			fmt.Println(user_id, "validated.")
			where_clause = fmt.Sprintf(` where user_id = '%s' `, user_id)
		}
	}

	query := fmt.Sprintf("SELECT id, user_id, subscriber_id, user_username, subscriber_name FROM user_subscriber_view%sORDER BY user_username, subscriber_name ASC LIMIT $1 OFFSET $2", where_clause)

	rows, err := d.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var user_subscriber_views []model.User_Subscriber_View
	for rows.Next() {
		var user_subscriber_view model.User_Subscriber_View
		if err := rows.Scan(&user_subscriber_view.Id, &user_subscriber_view.User_ID, &user_subscriber_view.Subscriber_Id, &user_subscriber_view.User_Username, &user_subscriber_view.Subscriber_Name); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_subscriber: %w", err)
		}

		user_subscriber_views = append(user_subscriber_views, user_subscriber_view)
	}
	return user_subscriber_views, nil
}

func (d *Database) UpdateUserSubscriber(ctx context.Context, user_subscriber model.User_Subscriber) error {
	fmt.Println("d UpdateUserSubscriber")

	query := `UPDATE user_subscriber SET user_id = $1, subscriber_id = $2 WHERE id = $3`
	fmt.Println("id", user_subscriber.Id)
	fmt.Println("user id", user_subscriber.User_ID)
	fmt.Println("subscriber id", user_subscriber.Subscriber_Id)

	_, err := d.DB.ExecContext(ctx, query, user_subscriber.User_ID, user_subscriber.Subscriber_Id, user_subscriber.Id)

	return err
}

func (d *Database) GetUserSubscriber(ctx context.Context, id string) (*model.User_Subscriber, error) {
	fmt.Println("d GetUserSubscriber")

	var user_subscriber model.User_Subscriber

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, user_id, subscriber_id FROM user_subscriber WHERE id = $1",
		id,
	).Scan(&user_subscriber.Id, &user_subscriber.Id, &user_subscriber.Id)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting subscriber: %w", err)
	}

	return &user_subscriber, nil
}

func (d *Database) CreateUserSubscriber(ctx context.Context, user_id string, subscriber_id string) (*model.User_Subscriber, error) {
	fmt.Println("d CreateUserSubscriber")

	user_subscriber := &model.User_Subscriber{
		Id:            uuid.New().String(),
		User_ID:       user_id,
		Subscriber_Id: subscriber_id,
	}

	query := `
        INSERT INTO user_subscriber (user_id, subscriber_id) VALUES ($1, $2)
    `

	_, err := d.DB.ExecContext(ctx, query,
		user_subscriber.User_ID,
		user_subscriber.Subscriber_Id,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user_subscriber: %w", err)
	}

	return user_subscriber, nil

}

func (d *Database) DeleteUserSubscriber(ctx context.Context, id string) error {
	fmt.Println("d DeleteUserSubscriber")

	query := `DELETE FROM user_subscriber WHERE id = $1`

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}

func (d *Database) LookupUserSubscriber(ctx context.Context, user_id string, subscriber_id string) (*model.User_Subscriber, error) {
	fmt.Println("d LookupUserSubscriber")

	var user_subscriber = model.User_Subscriber{}

	err := d.DB.QueryRowContext(ctx,
		"SELECT id, user_id, subscriber_id FROM user_subscriber WHERE user_id = $1 and subscriber_id = $2",
		user_id, subscriber_id,
	).Scan(&user_subscriber.Id, &user_subscriber.User_ID, &user_subscriber.Subscriber_Id)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user_subscriber: %w", err)
	}

	return &user_subscriber, nil
}

func (d *Database) LookupUserSubscribersByUserId(ctx context.Context, user_id string) ([]model.User_Subscriber_View, error) {
	fmt.Println("d LookupUserSubscribersByUser")

	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, user_id, subscriber_id, user_username, subscriber_name FROM user_subscriber_view WHERE user_id = $1",
		user_id,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing rows: %w", err)
	}
	defer rows.Close()

	var user_subscriber_views []model.User_Subscriber_View
	for rows.Next() {
		var user_subscriber_view model.User_Subscriber_View
		if err := rows.Scan(&user_subscriber_view.Id, &user_subscriber_view.User_ID, &user_subscriber_view.Subscriber_Id, &user_subscriber_view.User_Username, &user_subscriber_view.Subscriber_Name); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user_subscriber: %w", err)
		}

		user_subscriber_views = append(user_subscriber_views, user_subscriber_view)
	}
	return user_subscriber_views, nil
}
