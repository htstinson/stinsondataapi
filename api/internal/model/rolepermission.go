package model

import "time"

type Role_Permission struct {
	Role_Id       string    `json:"role_id"`
	Permission_Id string    `json:"permission_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type Role_Permission_View struct {
	Role_Id         string    `json:"role_id"`
	Role_Name       string    `json:"role_name"`
	Permission_Id   string    `json:"permission_id"`
	Permission_Name string    `json:"permission_name"`
	Object_Id       string    `json:"object_id"`
	Object_Name     string    `json:"object_name"`
	Object_Type     string    `json:"object_type"`
	CreatedAt       time.Time `json:"created_at"`
}
