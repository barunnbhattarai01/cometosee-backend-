package model

import "time"

type POSTDATA struct {
	ImageUrl  string    `firestore:"image_url"`
	Caption   string    `firestore:"caption"`
	Community string    `firestore:"community"`
	Username  string    `firestore:"username"`
	Created   time.Time `firestore:"created"`
}
