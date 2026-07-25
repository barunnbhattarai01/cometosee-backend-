package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type RequirementRepository interface {
	CreateRequirement(req model.Requirement) error
	UpdateRequirement(req model.Requirement) error
	DeleteRequirement(postID int) error
	GetRequirementByPost(postID int) (model.Requirement, error)
}

type requirementRepository struct{}

func NewRequirementRepository() RequirementRepository {
	return &requirementRepository{}
}

func (r *requirementRepository) CreateRequirement(req model.Requirement) error {

	query := `
	INSERT INTO post_requirements
	(
		post_id,
		min_age,
		max_age,
		gender,
		skill_level,
		verification_required,
		player_document_required,
		description
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := intailizer.DB.Exec(
		query,
		req.PostID,
		req.MinAge,
		req.MaxAge,
		req.Gender,
		req.SkillLevel,
		req.VerificationRequired,
		req.PlayerDocumentRequired,
		req.Description,
	)

	return err
}

func (r *requirementRepository) GetRequirementByPost(postID int) (model.Requirement, error) {

	var req model.Requirement

	query := `
	SELECT
	requirement_id,
	post_id,
	min_age,
	max_age,
	gender,
	skill_level,
	verification_required,
	player_document_required,
	description,
	created_at
	FROM post_requirements
	WHERE post_id=$1
	`

	err := intailizer.DB.QueryRow(query, postID).Scan(
		&req.RequirementID,
		&req.PostID,
		&req.MinAge,
		&req.MaxAge,
		&req.Gender,
		&req.SkillLevel,
		&req.VerificationRequired,
		&req.PlayerDocumentRequired,
		&req.Description,
		&req.CreatedAt,
	)

	return req, err
}

func (r *requirementRepository) UpdateRequirement(req model.Requirement) error {

	query := `
	UPDATE post_requirements
	SET
	min_age=$1,
	max_age=$2,
	gender=$3,
	skill_level=$4,
	verification_required=$5,
	player_document_required=$6,
	description=$7
	WHERE post_id=$8
	`

	_, err := intailizer.DB.Exec(
		query,
		req.MinAge,
		req.MaxAge,
		req.Gender,
		req.SkillLevel,
		req.VerificationRequired,
		req.PlayerDocumentRequired,
		req.Description,
		req.PostID,
	)

	return err
}

func (r *requirementRepository) DeleteRequirement(postID int) error {

	_, err := intailizer.DB.Exec(
		`DELETE FROM post_requirements WHERE post_id=$1`,
		postID,
	)

	return err
}
