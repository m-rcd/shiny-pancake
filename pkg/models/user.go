package models

type User struct {
	Id       string `json:"-"`
	Username string `json:"username"`
}
