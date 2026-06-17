package service

import (
	"cometosee/repository"
	"context"
)

type MapService struct {
	mapRepo *repository.MapRepository
}

func NewMapService(mapRepo *repository.MapRepository) *MapService {
	return &MapService{mapRepo: mapRepo}
}

func (s *MapService) GetMapEventPins(
	ctx context.Context,
	lat, lon float64,
	radius int,
	authid int,
) ([]map[string]interface{}, error) {

	if radius <= 0 {
		radius = 10000
	}

	sport, err := s.mapRepo.GetUserSport(authid)
	if err != nil {
		return nil, err
	}

	return s.mapRepo.MapEventPin(ctx, lat, lon, radius, sport)
}
