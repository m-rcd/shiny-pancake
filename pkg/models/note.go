package models

type Note struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	User    User   `json:"user"`
}
