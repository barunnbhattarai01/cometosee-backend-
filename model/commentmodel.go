package model

type Comment struct {
	ID      string      `json:"id"`
	UserID  string      `json:"user_id"`
	Comment string      `json:"comment"`
	Created interface{} `json:"created"`
}
