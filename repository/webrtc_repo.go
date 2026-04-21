package repository

import "cometosee/model"

type WebRTCRepository interface {
	SaveTrackInfo(trackInfo *model.TrackInfo) error
	RemoveTrackInfo(TrackId string) error
	GetTrackByRoomId(roomId string) ([]*model.TrackInfo, error)
	GetTrackByTrackId(trackId string) (*model.TrackInfo, error)
	AddPeertoRoom(roomId string, peerInfo *model.PeerInfo) error
	RemovePeerFromRoom(roomId string, peerId string) error
	GetPeersByRoomId(roomId string) ([]*model.PeerInfo, error)
	GetPeerByPeerId(peerId string) (*model.PeerInfo, error)
}

type webRTCRepo struct{}

func NewWebRTCRepository() WebRTCRepository {
	return &webRTCRepo{}
}
