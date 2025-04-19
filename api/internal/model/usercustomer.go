package model

import "time"

// User_Customrer
type User_Customer struct {
	Id           string    `json:"id"`
	User_ID      string    `json:"user_id"`
	Customer_Id  string    `json:"customer_id"`
	Assignedd_At time.Time `json:"assigned_at"`
}

// User_Customrer_View
type User_Customer_View struct {
	Id            string    `json:"id"`
	User_ID       string    `json:"user_id"`
	Customer_Id   string    `json:"customer_id"`
	User_Username string    `json:"user_username"`
	Customer_Name string    `json:"customer_name"`
	Assigned_At   time.Time `json:"assigned_at"`
}
