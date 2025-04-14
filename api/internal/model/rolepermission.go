package model

import "time"

type Role_Permission struct {
	Role_Id       string    `json:"role_id"`
	Permission_Id string    `json:"permission_id"`
	CreatedAt     time.Time `json:"created_at"`
}
