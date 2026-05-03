package model

import "time"

const (
	VideoCallSessionInitiated  = "initiated"
	VideoCallSessionInProgress = "in_progress"
	VideoCallSessionEnded      = "ended"
)

type VideoCallSession struct {
	Id                int64      `json:"id"`
	ConnectionID      int64      `json:"connection_id"`
	InitiatedByUserID int64      `json:"initiated_by_user_id"`
	AgoraChannelName  string     `json:"agora_channel_name"`
	Status            string     `json:"status"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	EndedAt           *time.Time `json:"ended_at,omitempty"`
	DurationSeconds   *int       `json:"duration_seconds,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
