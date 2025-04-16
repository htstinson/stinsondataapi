package model

import "time"

// Permission
type Permission struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Object_Id   string    `json:"object_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Permission_View struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Object_Id          string `json:"object_id"`
	Object_Name        string `json:"object_name"`
	Object_Description string `json:"object_description"`
	Object_Type        string `json:"object_type"`
}
