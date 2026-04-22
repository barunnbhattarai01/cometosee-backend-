package repository

import (
	"cometosee/model"
	"errors"
	"sync"
)

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

//peer contain peerconnection and datachannel
//peerconnection means audio/video streams ,ice connection ,not connection

//tracks means audio or video

type webRTCRepo struct {
	mu sync.RWMutex

	//tracks help to find the track's info using trackid
	tracks map[string]*model.TrackInfo // trackId to TrackInfo

	//roomsTrack help to find the groups of tracks using roomId
	roomsTrack map[string][]*model.TrackInfo // roomId to TrackInfo

	//peers help to find the specific user faster uding peerid
	peers map[string]*model.PeerInfo // peerId to PeerInfo

	// roomspeer help to find the users in specfic room and update the tracks
	roomspeer map[string]map[string]*model.PeerInfo // roomId to peerId to PeerInfo

}

func NewWebRTCRepository() WebRTCRepository {
	return &webRTCRepo{
		tracks:     make(map[string]*model.TrackInfo),
		roomsTrack: make(map[string][]*model.TrackInfo),

		peers:     make(map[string]*model.PeerInfo),
		roomspeer: make(map[string]map[string]*model.PeerInfo),
	}
}

func (r *webRTCRepo) SaveTrackInfo(trackInfo *model.TrackInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if trackInfo == nil || trackInfo.ID == "" {
		return errors.New("trackinfo is nil")
	}

	r.tracks[trackInfo.ID] = trackInfo

	return nil
}

func (r *webRTCRepo) RemoveTrackInfo(TrackId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tracks, TrackId)

	//also remove from roomtracks
	for roomid, tracks := range r.roomsTrack {
		var updated []*model.TrackInfo
		for _, t := range tracks {
			if t.ID != TrackId {
				updated = append(updated, t)
			}
		}
		r.roomsTrack[roomid] = updated
	}

	return nil
}

func (r *webRTCRepo) GetTrackByRoomId(roomId string) ([]*model.TrackInfo, error) {
	r.mu.RLock() //rlock used for readinfg operations
	defer r.mu.RUnlock()

	rooms, ok := r.roomsTrack[roomId]
	if !ok {
		return nil, errors.New("room not found")
	}

	var alltrack []*model.TrackInfo
	//loops over peers
	for _, tracks := range rooms {
		alltrack = append(alltrack, tracks)

	}
	return alltrack, nil
}

func (r *webRTCRepo) GetTrackByTrackId(trackId string) (*model.TrackInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	track, ok := r.tracks[trackId]
	if !ok {
		return nil, errors.New("track not found")
	}

	return track, nil
}

func (r *webRTCRepo) AddPeertoRoom(roomId string, peerInfo *model.PeerInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if peerInfo == nil || peerInfo.ID == "" {
		return errors.New("peerinfo is nil")
	}
	//save golabbly
	r.peers[peerInfo.ID] = peerInfo

	//intailize room if not exist
	if _, ok := r.roomspeer[roomId]; !ok {
		r.roomspeer[roomId] = make(map[string]*model.PeerInfo)
	}
	//add peer to room
	r.roomspeer[roomId][peerInfo.ID] = peerInfo

	return nil
}

func (r *webRTCRepo) RemovePeerFromRoom(roomId string, peerId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	//check room exist
	peersinroom, ok := r.roomspeer[roomId]
	if !ok {
		return errors.New("room not found")
	}
	//check peer exist in room
	if _, ok = peersinroom[peerId]; !ok {
		return errors.New("peer not found in room")
	}
	//remove peer from roo
	delete(peersinroom, peerId)
	//if room is empty then remove room
	if len(peersinroom) == 0 {
		delete(r.roomspeer, roomId)
	}

	//also remove from global peer list
	delete(r.peers, peerId)
	return nil
}

func (r *webRTCRepo) GetPeersByRoomId(roomId string) ([]*model.PeerInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peersinroom, ok := r.roomspeer[roomId]
	if !ok {
		return nil, errors.New("room not found")
	}

	var peers []*model.PeerInfo
	for _, peer := range peersinroom {
		peers = append(peers, peer)
	}

	return peers, nil
}

func (r *webRTCRepo) GetPeerByPeerId(peerId string) (*model.PeerInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peer, ok := r.peers[peerId]
	if !ok {
		return nil, errors.New("peer not found")
	}

	return peer, nil
}
