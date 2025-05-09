package model

import "time"

// User
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never send password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
	Roles        string    `json:"roles"`
	IP_address   string    `json:"ip_address"`
}

type CurrentUser struct {
	ID         string                  `json:"user_id"`
	Username   string                  `json:"username"`
	Roles      string                  `json:"roles"`
	IP_address string                  `json:"ip_address"`
	Subscribed []CurrentUserSubscribed `json:"subscribed"`
}

type CurrentUserSubscribed struct {
	ID            string `json:"id"`
	Subscriber_ID string `json:"subscriber_id"`
}
