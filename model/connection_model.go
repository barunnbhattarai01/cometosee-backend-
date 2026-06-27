package model

import "time"

type ConnectionStatus string

const (
	Pending  ConnectionStatus = "pending"
	Accepted ConnectionStatus = "accepted"
	Blocked  ConnectionStatus = "blocked"
)

type Connection struct {
	ID          string
	UserId1     string
	UserId2     string
	RequestedBy string
	Status      ConnectionStatus
	Avatar      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserPublic struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
