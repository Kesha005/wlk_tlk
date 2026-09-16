package stantion

import (
	"errors"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type OnePeer struct {
	Id string
	PC *webrtc.PeerConnection

	//This track which i receive from browser it will be forwarded to other peers
	RemoteTrack *webrtc.TrackRemote

	//This is local track which other peers, only speaker writes (half duplex)
	//While peer state is not speaking, his rtp packets not will write
	LocalTrack *webrtc.TrackLocalStaticRTP

	//Local track sender
	Sender *webrtc.RTPSender

	ws       *websocket.Conn
	mu       sync.Mutex
	Stantion *Stantion
	speaking bool
}

func NewRawPeer(id string) *OnePeer {
	return &OnePeer{
		Id:       id,
		speaking: false,
	}
}

func (p *OnePeer) OfferHandler(offer string) (*OnePeer, error) {

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return nil, err
	}
	audiotrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: webrtc.MimeTypeOpus,
	}, "pion", "wlk-tlk")

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return nil, err
	}

	sender, err := pc.AddTrack(audiotrack)

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return nil, err
	}
	p.PC = pc
	p.Sender = sender
	p.LocalTrack = p.LocalTrack

	err = pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer,
	})

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return nil, err
	}

	//TODO this will be handled differently it handle remote track on struct to access it from other places
	pc.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		if tr != nil {
			p.RemoteTrack = tr
		}

	})

	answer, err := pc.CreateAnswer(nil)

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return nil, err
	}

	err = pc.SetLocalDescription(answer)

	if err != nil {
		p.ws.WriteJSON(&WsMessage{
			Type: MSG_ERROR,
			Data: err.Error(),
		})

		return nil, err
	}

	//TODO handle closed conn write error
	p.ws.WriteJSON(WsMessage{
		Type: MSG_ANSWER,
		Data: pc.LocalDescription().SDP,
	})
	return p, err
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

// Peer state control will be used to log or ice status chekcking and controlling
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

	go func() {
		p.forwardLoop()
	}()

}

func (p *OnePeer) forwardLoop() {

	if p.RemoteTrack != nil {
		for {
			rtp, _, err := p.RemoteTrack.ReadRTP()
			if err != nil && !errors.Is(err, io.EOF) {
				p.ws.WriteJSON(&WsMessage{
					Type: MSG_ERROR,
					Data: err.Error(),
				})
				return
			}

			if !p.IsSpeaking() {
				continue
			}
			peers := p.Stantion.All()

			for _, peer := range peers {

				if peer.LocalTrack != nil {
					if p.Id == peer.Id {
						continue
					}
					peer.LocalTrack.WriteRTP(rtp)

				}
			}
		}

	}
}

func (p *OnePeer) Speak() {
	p.mu.Lock()
	p.speaking = true
	p.mu.Unlock()
	p.Stantion.SpeakStart(p.Id)
}

func (p *OnePeer) StartSpeak() {
	if !p.Stantion.CanSpeak(p.Id) {
		p.Speak()
	} else {
		return
	}
}

func (p *OnePeer) IsSpeaking() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.speaking
}
