package service

import "cometosee/repository"

type WebRTCService struct {
	repo repository.WebRTCRepository
}

func NewWebRTCService(repo repository.WebRTCRepository) *WebRTCService {
	return &WebRTCService{repo: repo}
}

func (s *WebRTCService) HandleWebRTCConnection() error {
	return s.repo.HandleWebRTCConnection()
}
