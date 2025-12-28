package model

import "time"

type Subscriber_Item_View struct {
	Id              string    `json:"id"`
	Item_ID         string    `json:"item_id"`
	Subscriber_Id   string    `json:"subscriber_id"`
	Item_Name       string    `json:"item_name"`
	Subscriber_Name string    `json:"subscriber_name"`
	CreatedAt_At    time.Time `json:"created_at"`
}
