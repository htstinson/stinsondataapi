package model

import "time"

// User_Customrer
type User_Customer_Roles struct {
	Id               string    `json:"id"`
	User_Customer_ID string    `json:"user_customer_id"`
	Role_Id          string    `json:"customer_id"`
	Created_At       time.Time `json:"assigned_at"`
	Updated_At       time.Time `json:"updated_at"`
}

type User_Customer_Roles_View struct {
	Id               string    `json:"id"`
	User_Customer_ID string    `json:"user_customer_id"`
	Role_Id          string    `json:"role_id"`
	Role_Name        string    `json:"role_name"`
	User_ID          string    `json:"user_id"`
	User_Name        string    `json:"user_username"`
	Customer_Id      string    `json:"customer_id"`
	Customer_Name    string    `json:"customer_name"`
	Created_At       time.Time `json:"created_at"`
	Updated_At       time.Time `json:"updated_at"`
}
