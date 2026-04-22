package service

import (
	"cometosee/model"
	"cometosee/repository"
	"errors"
)

type WebRTCService interface {
	AddPeer(roomId string, peer *model.PeerInfo) error
	RemovePeer(roomId, peerId string) error

	AddTrack(roomId string, track *model.TrackInfo) error
	RemoveTrack(trackId string) error

	GetRoomTracks(roomId string) ([]*model.TrackInfo, error)
	GetRoomPeers(roomId string) ([]*model.PeerInfo, error)
}

type webRTCService struct {
	repo repository.WebRTCRepository
}

func NewWebRTCService(repo repository.WebRTCRepository) WebRTCService {

	return &webRTCService{repo: repo}
}

func (s *webRTCService) AddPeer(roomId string, peer *model.PeerInfo) error {
	if roomId == "" || peer == nil || peer.ID == "" {
		return errors.New("invalid roomId or peer info")
	}

	return s.repo.AddPeertoRoom(roomId, peer)
}

func (s *webRTCService) RemovePeer(roomId, peerId string) error {
	return s.repo.RemovePeerFromRoom(roomId, peerId)
}

func (s *webRTCService) AddTrack(roomId string, track *model.TrackInfo) error {
	if roomId == "" || track == nil || track.ID == "" {
		return errors.New("invalid roomId or track info")
	}
	if err := s.repo.SaveTrackInfo(track); err != nil {
		return err
	}

	return nil
}

func (s *webRTCService) RemoveTrack(trackId string) error {
	return s.repo.RemoveTrackInfo(trackId)
}

func (s *webRTCService) GetRoomTracks(roomId string) ([]*model.TrackInfo, error) {
	return s.repo.GetTrackByRoomId(roomId)
}

func (s *webRTCService) GetRoomPeers(roomId string) ([]*model.PeerInfo, error) {
	return s.repo.GetPeersByRoomId(roomId)
}
