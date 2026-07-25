package model

import "time"

type Requirement struct {
	RequirementID          int       `json:"requirement_id"`
	PostID                 int       `json:"post_id"`
	MinAge                 *int      `json:"min_age"`
	MaxAge                 *int      `json:"max_age"`
	Gender                 string    `json:"gender"`
	SkillLevel             string    `json:"skill_level"`
	VerificationRequired   bool      `json:"verification_required"`
	PlayerDocumentRequired bool      `json:"player_document_required"`
	Description            string    `json:"description"`
	CreatedAt              time.Time `json:"created_at"`
}
