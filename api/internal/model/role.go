package model

type Roles struct {
	Id       string
	Username string
	Names    string
}

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
