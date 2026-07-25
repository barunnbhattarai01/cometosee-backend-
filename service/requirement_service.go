package service

import (
	"cometosee/model"
	"cometosee/repository"
)

type RequirementService interface {
	CreateRequirement(req model.Requirement) error
	UpdateRequirement(req model.Requirement) error
	DeleteRequirement(postID int) error
	GetRequirement(postID int) (model.Requirement, error)
}

type requirementService struct {
	repo repository.RequirementRepository
}

func NewRequirementService(repo repository.RequirementRepository) RequirementService {
	return &requirementService{repo: repo}
}

func (s *requirementService) CreateRequirement(req model.Requirement) error {
	return s.repo.CreateRequirement(req)
}

func (s *requirementService) UpdateRequirement(req model.Requirement) error {
	return s.repo.UpdateRequirement(req)
}

func (s *requirementService) DeleteRequirement(postID int) error {
	return s.repo.DeleteRequirement(postID)
}

func (s *requirementService) GetRequirement(postID int) (model.Requirement, error) {
	return s.repo.GetRequirementByPost(postID)
}
