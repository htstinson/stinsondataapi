package model

import "time"

// User_Customrer
type User_Subscriber struct {
	Id            string    `json:"id"`
	User_ID       string    `json:"user_id"`
	Subscriber_Id string    `json:"subscriber_id"`
	Assigned_At   time.Time `json:"assigned_at"`
}

// User_Customrer_View
type User_Subscriber_View struct {
	Id              string    `json:"id"`
	User_ID         string    `json:"user_id"`
	Subscriber_Id   string    `json:"subscriber_id"`
	User_Username   string    `json:"user_username"`
	Subscriber_Name string    `json:"subscriber_name"`
	Assigned_At     time.Time `json:"assigned_at"`
}
