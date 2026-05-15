package model

import "time"

type POSTDATA struct {
	ImageUrl  string    `json:"image_url"`
	Caption   string    `json:"caption"`
	Community string    `json:"community"`
	Username  string    `json:"username"`
	Created   time.Time `json:"created_at"`
}

type PostSlot struct {
	SlotID          int       `json:"slot_id"`
	PostID          int       `json:"post_id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	MaxParticipants int       `json:"max_participants"`
	CreatedAt       time.Time `json:"created_at"`
}

type SlotParticipant struct {
	ID       int       `json:"id"`
	SlotID   int       `json:"slot_id"`
	AuthID   int       `json:"auth_id"`
	JoinedAt time.Time `json:"joined_at"`
}
