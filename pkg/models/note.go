package models

type Note struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	User     User   `json:"user"`
	Archived bool   `json:"archived"`
}
