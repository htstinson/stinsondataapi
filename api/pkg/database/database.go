package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/htstinson/stinsondataapi/api/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	// Item Get, Create, List, Update, Delete
	GetItem(ctx context.Context, id string) (*model.Item, error)
	CreateItem(ctx context.Context, item *model.Item) error
	ListItems(ctx context.Context, limit, offset int) ([]model.Item, error)
	UpdateItem(cts context.Context, item *model.Item) error
	DeleteItem(ctx context.Context, id string) error

	// User Get(2), Create, List, Update, Delete
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, username string, password string) (*model.User, error)
	SelectUsers(ctx context.Context, limit, offset int) ([]model.User, error)
	UpdateUser(cts context.Context, item *model.User) error
	DeleteUser(ctx context.Context, id string) error

	// Blocked
	SelectBlocked(ctx context.Context, limit, offset int, sort string, order string) ([]model.Blocked, error)
	GetBlocked(ctx context.Context, id string) (*model.Blocked, error)
	UpdateBlocked(cts context.Context, item *model.Blocked) error
	CreateBlocked(ctx context.Context, blocked model.Blocked) (*model.Blocked, error)
	DeleteBlocked(ctx context.Context, id string) error

	RowCount(tablename string) (int, error)

	SelectRoles(ctx context.Context, userID string) (model.Roles, error)
	SelectUserRoles(ctx context.Context, limit, offset int) ([]model.User, error)

	Close() error
}

type Database struct {
	db *sql.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func New(cfg Config) (Repository, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s\n",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Initialize schema
	if err := initializeSchema(db); err != nil {
		return nil, fmt.Errorf("error initializing schema: %w", err)
	}

	return &Database{db: db}, nil
}

func initializeSchema(db *sql.DB) error {

	// Create Users table
	queries := []string{
		`CREATE TABLE IF NOT EXISTS blocked (
            id VARCHAR(36) PRIMARY KEY,
            ip VARCHAR(20) UNIQUE NOT NULL,
            notes VARCHAR(255),
            created_at TIMESTAMP WITH TIME ZONE NOT NULL
        )`,
		`ALTER TABLE blocked ALTER COLUMN id SET DEFAULT gen_random_uuid();`,
		`ALTER TABLE blocked ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;`,
		`CREATE INDEX IF NOT EXISTS blocked_ip_idx ON blocked(ip)`,
		`CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(36) PRIMARY KEY,
            username VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL
        )`,
		`CREATE INDEX IF NOT EXISTS users_username_idx ON users(username)`,
		`CREATE TABLE IF NOT EXISTS people (
            id VARCHAR(36) PRIMARY KEY,
            firstname VARCHAR(255) UNIQUE NOT NULL,
            lastname VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL
        )`,
		`CREATE INDEX IF NOT EXISTS users_username_idx ON users(username)`,
		`CREATE TABLE IF NOT EXISTS items (
            id VARCHAR(36) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL
        );
        CREATE INDEX IF NOT EXISTS items_created_at_idx ON items(created_at DESC);`,
		`CREATE TABLE IF NOT EXISTS accounts (
            id VARCHAR(36) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL,
			phone VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL
        );
        CREATE INDEX IF NOT EXISTS accounts_created_at_idx ON items(created_at DESC);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error creating schema: %w", err)
		}
	}
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// Item

func (d *Database) GetItem(ctx context.Context, id string) (*model.Item, error) {
	var item model.Item
	err := d.db.QueryRowContext(ctx,
		"SELECT id, name, created_at FROM items WHERE id = $1",
		id,
	).Scan(&item.ID, &item.Name, &item.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting item: %w", err)
	}
	return &item, nil
}

func (d *Database) UpdateItem(ctx context.Context, item *model.Item) error {

	query := `UPDATE items SET name = $1 WHERE id = $2`

	_, err := d.db.ExecContext(ctx, query, item.Name, item.ID)

	return err

}

func (d *Database) DeleteItem(ctx context.Context, id string) error {

	query := `DELETE FROM items WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}

func (d *Database) CreateItem(ctx context.Context, item *model.Item) error {
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()

	_, err := d.db.ExecContext(ctx,
		"INSERT INTO items (id, name, created_at) VALUES ($1, $2, $3)",
		item.ID, item.Name, item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating item: %w", err)
	}
	return nil
}

func (d *Database) ListItems(ctx context.Context, limit, offset int) ([]model.Item, error) {
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, name, created_at FROM items ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error listing items: %w", err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

// any table
func (d *Database) RowCount(tablename string) (int, error) {
	fmt.Println("d RowCount")

	var count int

	q := fmt.Sprintf("SELECT COUNT(*) FROM %s", tablename)
	fmt.Println(q)

	ctx := context.Background()

	rows, err := d.db.QueryContext(ctx, q)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return count, err

}

// Admin - Blocked
func (d *Database) SelectBlocked(ctx context.Context, limit int, offset int, sort string, order string) ([]model.Blocked, error) {

	q := fmt.Sprintf("SELECT id, ip, notes, created_at FROM blocked ORDER BY %s %s LIMIT $1 OFFSET $2", sort, strings.ToUpper(order))

	rows, err := d.db.QueryContext(ctx, q, limit, offset)

	if err != nil {
		fmt.Printf("[%v] [database][SelectBlocked] error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
		fmt.Printf("[%v] [database][SelectBlocked] query: %s.\n", time.Now().Format(time.RFC3339), q)
		fmt.Printf("[%v] [database][SelectBlocked] limit %v offset %v sort %s order %s.\n", time.Now().Format(time.RFC3339), limit, offset, strings.ToUpper(sort), strings.ToUpper(order))
		return nil, fmt.Errorf("query error")
	}
	defer rows.Close()

	var items []model.Blocked

	for rows.Next() {
		var item model.Blocked
		var notesNullable sql.NullString
		if err := rows.Scan(&item.ID, &item.IP, &notesNullable, &item.CreatedAt); err != nil {
			fmt.Printf("[%v] [database][ListBlocked] error: %s.\n", time.Now().Format(time.RFC3339), err.Error())
			return nil, fmt.Errorf("error scanning blocked: %w", err)
		}
		if notesNullable.Valid {
			item.Notes = notesNullable.String
		} else {
			item.Notes = ""
		}

		items = append(items, item)
	}

	return items, nil
}

func (d *Database) UpdateBlocked(ctx context.Context, blocked *model.Blocked) error {
	fmt.Println("d UpdateBlocked", blocked.IP, blocked.Notes)

	query := `UPDATE blocked SET ip=$1, notes=$2 WHERE id = $3`

	_, err := d.db.ExecContext(ctx, query, blocked.IP, blocked.Notes, blocked.ID)

	return err

}

func (d *Database) GetBlockedByIP(ctx context.Context, ip string) (*model.Blocked, error) {
	var blocked model.Blocked
	var notesNull sql.NullString // Use sql.NullString to handle NULL values

	query := fmt.Sprintf(`SELECT id, ip, notes, created_at FROM blocked WHERE ip ='%s'`, ip)

	err := d.db.QueryRowContext(ctx, query).Scan(&blocked.ID, &blocked.IP, &notesNull, &blocked.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		fmt.Println("d GetBlockedByIP", err.Error())
		return nil, fmt.Errorf("error getting blocked: %w", err)
	}

	// Only assign the value if it's not NULL
	if notesNull.Valid {
		blocked.Notes = notesNull.String
	} else {
		blocked.Notes = "" // Or another default value of your choice
	}

	return &blocked, nil
}

func (d *Database) GetBlocked(ctx context.Context, id string) (*model.Blocked, error) {
	var blocked model.Blocked
	var notesNull sql.NullString // Use sql.NullString to handle NULL values

	query := "SELECT id, ip, notes, created_at FROM blocked WHERE id = $1"

	err := d.db.QueryRowContext(ctx, query, id).Scan(&blocked.ID, &blocked.IP, &notesNull, &blocked.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		fmt.Println("d GetBlocked", err.Error())
		return nil, fmt.Errorf("error getting blocked: %w", err)
	}

	// Only assign the value if it's not NULL
	if notesNull.Valid {
		blocked.Notes = notesNull.String
	} else {
		blocked.Notes = "" // Or another default value of your choice
	}

	return &blocked, nil
}

func (d *Database) CreateBlocked(ctx context.Context, blocked model.Blocked) (*model.Blocked, error) {
	fmt.Println("d CreateBlocked(ctx, blocked)")
	blocked.CreatedAt = time.Now()

	_, err := d.GetBlockedByIP(ctx, blocked.IP)
	if err == nil {
		return nil, errors.New("duplicate")
	}

	query := `
        INSERT INTO blocked (ip, notes, created_at)
        VALUES ($1, $2, $3)
    `

	_, err = d.db.ExecContext(ctx, query, blocked.IP, blocked.Notes, blocked.CreatedAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error creating blocked: %w", err)
	}

	return &blocked, nil

}

func (d *Database) DeleteBlocked(ctx context.Context, id string) error {

	query := `DELETE FROM blocked WHERE id = $1`

	_, err := d.db.ExecContext(ctx, query, id)

	return err
}

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
		"SELECT id, username, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)

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

//Roles

func (d *Database) SelectRoles(ctx context.Context, userId string) (model.Roles, error) {
	fmt.Println("d SelectRoles")

	var roles = model.Roles{}

	query := `
        SELECT user_id, username, role_name
        FROM user_roles_view
        WHERE user_id = $1
    `

	err := d.db.QueryRowContext(ctx, query, userId).Scan(
		&roles.Id, &roles.Username, &roles.Names,
	)

	return roles, err
}
