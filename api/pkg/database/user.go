package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/htstinson/stinsondataapi/api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

//User

func (d *Database) SelectUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	fmt.Println("database.go SelectUsers()")
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, username, ip_address FROM users ORDER BY username ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.IP_address); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		user.Roles = strings.Replace(user.Roles, "{", "", -1)
		user.Roles = strings.Replace(user.Roles, "}", "", -1)
		users = append(users, user)
	}
	return users, nil
}

func (d *Database) SelectUserRoles(ctx context.Context, limit, offset int) ([]model.User, error) {
	fmt.Println("database.go SelectUserRoles()")
	rows, err := d.db.QueryContext(ctx,
		"SELECT user_id, username, ip_address, role_name FROM user_roles_view ORDER BY username ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.IP_address, &user.Roles); err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		user.Roles = strings.Replace(user.Roles, "{", "", -1)
		user.Roles = strings.Replace(user.Roles, "}", "", -1)
		users = append(users, user)
	}
	return users, nil
}

func (d *Database) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := d.db.QueryRowContext(ctx,
		"SELECT id, username, ipaddress, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.IP_address, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (d *Database) CreateUser(ctx context.Context, username, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	query := `
        INSERT INTO users (id, username, password_hash, created_at)
        VALUES ($1, $2, $3, $4)
    `

	_, err = d.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (d *Database) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	fmt.Printf("[%v] [GetUserByUsername] %s\n", time.Now().Format(time.RFC3339), username)

	user := &model.User{}
	query := `
        SELECT id, username, password_hash, created_at
        FROM users
        WHERE username = $1
    `

	err := d.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (d *Database) DeleteUser(ctx context.Context, id string) error {

	query := `DELETE FROM users WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}

func (d *Database) UpdateUser(ctx context.Context, user *model.User) error {
	fmt.Println("h UpdateUser")

	query := `UPDATE users SET username = $1, ip_address = $2 WHERE id = $3`

	fmt.Println("user.Ip_address", user.IP_address)

	var ipAddress string
	if user.IP_address == "" {
		ipAddress = "0.0.0.0" // or another default IP
	} else {
		ipAddress = user.IP_address
	}

	fmt.Println(ipAddress)

	_, err := d.db.ExecContext(ctx, query, user.Username, ipAddress, user.ID)

	return err
}
