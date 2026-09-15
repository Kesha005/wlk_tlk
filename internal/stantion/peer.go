package stantion

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type OnePeer struct {
	Id         string
	PC         *webrtc.PeerConnection
	AudioTrack *webrtc.TrackLocalStaticRTP
	Sender     *webrtc.RTPSender
	ws         *websocket.Conn
	mu         sync.Mutex
	speaking   bool
}

func NewRawPeer(id string) *OnePeer {
	return &OnePeer{
		Id:       id,
		speaking: false,
	}
}

func (p *OnePeer) OfferHandler(offer string) {

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return
	}
	audiotrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: webrtc.MimeTypeOpus,
	}, "pion", "wlk-tlk")

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return
	}

	sender, err := pc.AddTrack(audiotrack)

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return
	}
	p.PC = pc
	p.Sender = sender
	p.AudioTrack = audiotrack

	err = pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer,
	})

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return
	}

	sender.Stop()

	answer, err := pc.CreateAnswer(nil)

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return
	}

	err = pc.SetLocalDescription(answer)

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return
	}

	//TODO handle closed conn write error
	p.ws.WriteJSON(WsMessage{
		Type: MSG_ANSWER,
		Data: pc.LocalDescription().SDP,
	})
}

func (p *OnePeer) IceHandler(ice *webrtc.ICECandidateInit) {

	if ice != nil {
		if p.PC != nil {
			err := p.PC.AddICECandidate(*ice)

			if err != nil {
				p.ws.WriteJSON(&WsMessage{
					Type: MSG_ERROR,
					Data: err.Error(),
				})

			}
		}
	}
}

func (p *OnePeer) RunStateControl() {

	p.PC.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i != nil {

			if p.ws != nil {
				ice := i.ToJSON()
				err := p.ws.WriteJSON(&WsMessage{
					Type: MSG_ICE,
					Data: "",
					Ice:  &ice,
				})

				if err != nil {
					p.ws.WriteJSON(&WsMessage{
						Type: MSG_ERROR,
						Data: err.Error(),
					})

				}

			}
		}
	})

}
