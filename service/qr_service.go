package service

import (
	"cometosee/model"
	"cometosee/repository"
)

type QRService struct {
	repo *repository.QRRepository
}

func NewQRService(repo *repository.QRRepository) *QRService {
	return &QRService{
		repo: repo,
	}
}

func (s *QRService) GetJoinedEvents(authID int) ([]model.JoinedEventQR, error) {
	return s.repo.GetJoinedEvents(authID)
}

func (s *QRService) GetOwnerPosts(authID int) ([]model.OwnerPostQR, error) {
	return s.repo.GetOwnerPosts(authID)
}

func (s *QRService) VerifyQR(req model.VerifyParticipantRequest, ownerID int) (*model.QRParticipant, error) {

	if err := s.repo.VerifyQR(req.Token, ownerID); err != nil {
		return nil, err
	}

	return s.repo.GetParticipant(req.Token)
}
