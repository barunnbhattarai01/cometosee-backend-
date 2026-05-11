package model

import "time"

type POSTDATA struct {
	ImageUrl  string    `json:"image_url"`
	Caption   string    `json:"caption"`
	Community string    `json:"community"`
	Username  string    `json:"username"`
	Created   time.Time `json:"created_at"`
}
