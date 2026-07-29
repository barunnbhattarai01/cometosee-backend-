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
	ID     int `json:"id"`
	SlotID int `json:"slot_id"`
	AuthID int `json:"auth_id"`

	QRToken     string     `json:"qr_token"`
	QRExpiresAt *time.Time `json:"qr_expires_at,omitempty"`

	CheckedIn   bool       `json:"checked_in"`
	CheckedInAt *time.Time `json:"checked_in_at,omitempty"`
	CheckedInBy *int       `json:"checked_in_by,omitempty"`

	JoinedAt time.Time `json:"joined_at"`
}

type CancelPostInfo struct {
	Caption   string
	Venue     string
	Sport     string
	StartTime time.Time
	EndTime   time.Time
}

type VerifyParticipantRequest struct {
	Token string `json:"token"`
}

type ParticipantQRResponse struct {
	Token string `json:"token"`
}

type JoinedEventQR struct {
	PostID    int    `json:"post_id"`
	SlotID    int    `json:"slot_id"`
	Caption   string `json:"caption"`
	Venue     string `json:"venue"`
	QRToken   string `json:"qr_token"`
	CheckedIn bool   `json:"checked_in"`
}

type OwnerPostQR struct {
	PostID  int    `json:"post_id"`
	Caption string `json:"caption"`
	Venue   string `json:"venue"`
}

type QRParticipant struct {
	AuthID      int     `json:"auth_id"`
	Username    string  `json:"username"`
	CallingName *string `json:"calling_name,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Venue       *string `json:"venue"`
	CheckedIn   bool    `json:"checked_in"`

	PostID      int    `json:"post_id"`
	PostCaption string `json:"post_caption"`
}
