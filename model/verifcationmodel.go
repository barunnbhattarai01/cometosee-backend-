package model

import "time"

type Verification struct {
	VerificationID   int        `json:"verification_id"`
	AuthID           int        `json:"auth_id"`
	CitizenshipFront string     `json:"citizenship_front"`
	CitizenshipBack  string     `json:"citizenship_back"`
	Status           string     `json:"status"`
	RejectionReason  *string    `json:"rejection_reason"`
	VerifiedAt       *time.Time `json:"verified_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type PlayerDocument struct {
	DocumentID      int        `json:"document_id"`
	AuthID          int        `json:"auth_id"`
	DocumentName    string     `json:"document_name"`
	DocumentType    string     `json:"document_type"`
	DocumentURL     string     `json:"document_url"`
	IssuedBy        string     `json:"issued_by"`
	IssueDate       string     `json:"issue_date"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason"`
	VerifiedAt      *time.Time `json:"verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
}
