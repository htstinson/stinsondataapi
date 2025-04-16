package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/htstinson/stinsondataapi/api/internal/model"

	_ "github.com/lib/pq"
)

type Repository interface {
	// Item
	GetItem(ctx context.Context, id string) (*model.Item, error)
	CreateItem(ctx context.Context, item *model.Item) error
	SelectItems(ctx context.Context, limit, offset int) ([]model.Item, error)
	UpdateItem(ctx context.Context, item *model.Item) error
	DeleteItem(ctx context.Context, id string) error

	// User
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, username string, password string) (*model.User, error)
	SelectUsers(ctx context.Context, limit, offset int) ([]model.User, error)
	UpdateUser(ctx context.Context, item *model.User) error
	DeleteUser(ctx context.Context, id string) error

	// Customer
	GetCustomer(ctx context.Context, id string) (*model.Customer, error)
	GetCustomerByName(ctx context.Context, name string) (*model.Customer, error)
	CreateCustomer(ctx context.Context, name string) (*model.Customer, error)
	SelectCustomers(ctx context.Context, limit, offset int) ([]model.Customer, error)
	UpdateCustomer(ctx context.Context, customer *model.Customer) error
	DeleteCustomer(ctx context.Context, id string) error
	SelectUserRoles(ctx context.Context, limit, offset int) ([]model.User, error)

	// Blocked
	SelectBlocked(ctx context.Context, limit, offset int, sort string, order string) ([]model.Blocked, error)
	GetBlocked(ctx context.Context, id string) (*model.Blocked, error)
	UpdateBlocked(ctx context.Context, item *model.Blocked) error
	CreateBlocked(ctx context.Context, blocked model.Blocked) (*model.Blocked, error)
	DeleteBlocked(ctx context.Context, id string) error

	// Roles
	SelectRolesByUser(ctx context.Context, userID string) (model.Roles, error)
	SelectRoles(ctx context.Context, limit, offset int) ([]model.Role, error)
	GetRole(ctx context.Context, id string) (*model.Role, error)
	UpdateRole(ctx context.Context, role *model.Role) error
	CreateRole(ctx context.Context, name string) (*model.Role, error)
	DeleteRole(ctx context.Context, id string) error

	// Permission
	GetPermission(ctx context.Context, id string) (*model.Permission, error)
	CreatePermission(ctx context.Context, name string, description string) (*model.Permission, error)
	SelectPermissions(ctx context.Context, limit, offset int) ([]model.Permission, error)
	SelectPermissions_View(ctx context.Context, limit, offset int) ([]model.Permission_View, error)
	UpdatePermission(ctx context.Context, permission *model.Permission) error
	DeletePermission(ctx context.Context, id string) error

	// User Permissions
	SelectUserPermissions(ctx context.Context, limit, offset int) ([]model.User_Permission, error)

	//Role Permissions
	SelectRolePermissions(ctx context.Context, limit, offset int) ([]model.Role_Permission_View, error)

	RowCount(tablename string) (int, error)

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
