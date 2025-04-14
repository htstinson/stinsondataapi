package model

import "time"

type User_Permission struct {
	User_Id       string    `json:"user_id"`
	Permission_Id string    `json:"permission_id"`
	CreatedAt     time.Time `json:"created_at"`
}
