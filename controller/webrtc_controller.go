package controller

import "cometosee/service"

//so udp(user datagram protocol) doesn't require a connection to be established before data is sent, it is often used for real-time applications like video conferencing, online gaming, and live streaming where low latency is crucial.
//  In the context of WebRTC, UDP is commonly used for transmitting media streams (audio and video) between peers, allowing for faster communication compared to TCP (Transmission Control Protocol) which requires a connection to be established and has more overhead for error checking and retransmission.
//In WebRTC, the Real-time Transport Protocol (RTP) is typically used over UDP to transmit media streams. RTP provides mechanisms for sequencing, timestamping, and payload identification, which are essential for synchronizing audio and video streams in real-time communication.
//  Additionally, WebRTC can also use TCP as a fallback option if UDP is not available or if there are network restrictions that prevent UDP traffic.
//so first webrtc try through stun server and if it fails then it will try through turn server.
//Nats convert private/local ip to public ip used by router

type WebRTCController struct {
	service *service.WebRTCService
}

func NewWebRTCController(service *service.WebRTCService) *WebRTCController {
	return &WebRTCController{service: service}
}
