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
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Object_Id            string `json:"object_id"`
	V_Object_Name        string `json:"v_object_name"`
	V_Object_Description string `json:"v_object_description"`
	V_Object_Type        string `json:"v_object_type"`
}
