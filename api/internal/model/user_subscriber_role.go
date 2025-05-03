package model

import "time"

// User_Subscriber_Role
type User_Subscriber_Role struct {
	Id                 string    `json:"id"`
	User_Subscriber_ID string    `json:"user_subscriber_id"`
	Role_Id            string    `json:"role_id"`
	Created_At         time.Time `json:"assigned_at"`
	Updated_At         time.Time `json:"updated_at"`
}

type User_Subscriber_Roles_View struct {
	Id                 string    `json:"id"`
	User_Subscriber_ID string    `json:"user_subscriber_id"`
	Role_Id            string    `json:"role_id"`
	Role_Name          string    `json:"role_name"`
	User_ID            string    `json:"user_id"`
	User_Name          string    `json:"user_username"`
	Subscriber_Id      string    `json:"subscriber_id"`
	Subscriber_Name    string    `json:"subscriber_name"`
	Created_At         time.Time `json:"created_at"`
	Updated_At         time.Time `json:"updated_at"`
}
