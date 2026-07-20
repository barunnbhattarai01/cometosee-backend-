package service

import (
	"cometosee/model"
	"cometosee/repository"
)

type DiscoveryService interface {
	DiscoverUsers(
		lat float64,
		lon float64,
		radius int,
		sport string,
		skill string,
		currentUserId int,
	) ([]model.UserDetailInfoForCreate, error)
}

type discoveryService struct {
	repo repository.InterestAndLocationDiscoveryRepository
}

func NewDiscoveryService(r repository.InterestAndLocationDiscoveryRepository) DiscoveryService {
	return &discoveryService{repo: r}
}

func (s *discoveryService) DiscoverUsers(
	lat float64,
	lon float64,
	radius int,
	sport string,
	skill string,
	currentUserId int,
) ([]model.UserDetailInfoForCreate, error) {

	users, err := s.repo.FindNearbyUsers(lat, lon, radius, sport, skill, currentUserId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
