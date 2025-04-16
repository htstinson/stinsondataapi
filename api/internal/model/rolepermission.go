package model

import "time"

type Role_Permission struct {
	Role_Id       string    `json:"role_id"`
	Permission_Id string    `json:"permission_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type Role_Permission_View struct {
	Role_Id           string    `json:"role_id"`
	V_Role_Name       string    `json:"v_role_name"`
	Permission_Id     string    `json:"permission_id"`
	V_Permission_Name string    `json:"v_permission_name"`
	Object_Id         string    `json:"object_id"`
	V_Object_Name     string    `json:"v_object_name"`
	V_Object_Type     string    `json:"v_object_type"`
	CreatedAt         time.Time `json:"created_at"`
}
