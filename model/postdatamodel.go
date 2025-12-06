package model

import "time"

type POSTDATA struct {
	ImageUrl  string    `firestore:"image_url"`
	Caption   string    `firestore:"caption"`
	Community string    `firestore:"community"`
	Created   time.Time `firestore:"created"`
}
