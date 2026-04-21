package model

import "github.com/pion/webrtc/v3"

type PeerInfo struct {
	ID             string
	PeerConnection *webrtc.PeerConnection //main  webrtc connection
	DataChannel    *webrtc.DataChannel    //message channel for signaling and other data exchange
}

type TrackInfo struct {
	ID    string
	Kind  string // audio or video
	Track *webrtc.TrackLocalStaticSample
	SSRC  uint32 //ssrc is unique id for each track in webrtc session
}

type SignalingMessage struct {
	Type      string `json:"type"` // offer, answer, candidate
	SDP       string `json:"sdp,omitempty"`
	Candidate string `json:"candidate,omitempty"`
}
