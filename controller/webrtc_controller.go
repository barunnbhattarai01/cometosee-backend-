package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//so udp(user datagram protocol) doesn't require a connection to be established before data is sent, it is often used for real-time applications like video conferencing, online gaming, and live streaming where low latency is crucial.
//  In the context of WebRTC, UDP is commonly used for transmitting media streams (audio and video) between peers, allowing for faster communication compared to TCP (Transmission Control Protocol) which requires a connection to be established and has more overhead for error checking and retransmission.
//In WebRTC, the Real-time Transport Protocol (RTP) is typically used over UDP to transmit media streams. RTP provides mechanisms for sequencing, timestamping, and payload identification, which are essential for synchronizing audio and video streams in real-time communication.
//  Additionally, WebRTC can also use TCP as a fallback option if UDP is not available or if there are network restrictions that prevent UDP traffic.
//so first webrtc try through stun server and if it fails then it will try through turn server.
//Nats convert private/local ip to public ip used by router

type WebRTCController struct {
	service service.WebRTCService
}

func NewWebRTCController(service service.WebRTCService) *WebRTCController {
	return &WebRTCController{service: service}
}

// join room
func (c *WebRTCController) JoinRoom(w http.ResponseWriter, r *http.Request) {
	type JoinRoomRequest struct {
		RoomID string `json:"roomId"`
		PeerID string `json:"peerId"`
	}

	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := c.service.AddPeer(req.RoomID, &model.PeerInfo{ID: req.PeerID})
	if err != nil {
		http.Error(w, "Failed to join room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	common.WriteJSONMessage(w, "sucessfully  join room")

}

// leave room
func (c *WebRTCController) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	type LeaveRoomRequest struct {
		RoomID string `json:"roomId"`
		PeerID string `json:"peerId"`
	}
	var req LeaveRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := c.service.RemovePeer(req.RoomID, req.PeerID)
	if err != nil {
		http.Error(w, "Failed to leave room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	common.WriteJSONMessage(w, "sucessfully  leave room")
}

// get peers in room
func (c *WebRTCController) GetPeersInRoom(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	roomId := id["roomId"]

	if roomId == "" {
		http.Error(w, "Missing roomId parameter", http.StatusBadRequest)
		return
	}
	peers, err := c.service.GetRoomPeers(roomId)
	if err != nil {
		http.Error(w, "Failed to get peers in room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sucessfully fetched",
		"data":    peers,
	})
}

// add track
func (c *WebRTCController) AddTrack(w http.ResponseWriter, r *http.Request) {
	type AddTrackRequest struct {
		RoomID  string `json:"roomId"`
		TrackID string `json:"trackId"`
		Kind    string `json:"kind"`
	}
	var req AddTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	track := &model.TrackInfo{
		ID:     req.TrackID,
		RoomID: req.RoomID,
		Kind:   req.Kind,
	}

	err := c.service.AddTrack(req.RoomID, track)
	if err != nil {
		http.Error(w, "Failed to add track: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	common.WriteJSONMessage(w, "sucessfully  add track")
}

// remove track
func (c *WebRTCController) RemoveTrack(w http.ResponseWriter, r *http.Request) {
	type RemoveTrackRequest struct {
		TrackID string `json:"trackId"`
	}
	var req RemoveTrackRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := c.service.RemoveTrack(req.TrackID)
	if err != nil {
		http.Error(w, "Failed to remove track: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	common.WriteJSONMessage(w, "sucessfully  remove track")
}

//get tracks in room

func (c *WebRTCController) GetTracksInRoom(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)
	roomId := id["roomId"]

	if roomId == "" {
		http.Error(w, "Missing roomId parameter", http.StatusBadRequest)
		return
	}
	tracks, err := c.service.GetRoomTracks(roomId)
	if err != nil {
		http.Error(w, "Failed to get tracks in room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sucessfully fetched",
		"data":    tracks,
	})
}
