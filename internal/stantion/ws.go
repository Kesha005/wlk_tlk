package stantion

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

const (
	MSG_SDP = "sdp"
	MSG_ICE = "ice"

	MSG_ANSWER = "answer"

	MSG_ERROR = "error"

	MSG_SPEAK_START = "s_start"

	MSG_SPEAK_END = "s_end"
)

type WsMessage struct {
	Type string                   `json:"type"`
	Data string                   `json:"data"`
	Ice  *webrtc.ICECandidateInit `json:"ice"`
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 2048,
	ReadBufferSize:  2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (st *Stantion) wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Failed to make websocket connection", 500)
	}

	id := uuid.New().String()
	peer := NewRawPeer(id)
	peer.ws = conn

	st.InitPeer(peer)

	defer conn.Close()

}

func (st *Stantion) InitPeer(peer *OnePeer) {
	go func() {
		st.peerWsMsgHandler(peer)
	}()

}

func (st *Stantion) peerWsMsgHandler(peer *OnePeer) {
	var msg WsMessage
	for {
		err := peer.ws.ReadJSON(&msg)

		if err != nil {
			//TODO handle error there
		}
		switch msg.Type {
		case MSG_SDP:
			//Add peer to stantion
			newPeer, err := peer.OfferHandler(msg.Data)
			if err != nil {
				st.Add(newPeer)
			}

		case MSG_ICE:
			peer.IceHandler(msg.Ice)
		case MSG_ERROR:
			fmt.Println(msg.Data)
			return
		case MSG_SPEAK_END:
			//TODO handle spek
		case MSG_SPEAK_START:
			//TODO handle stop

		}

	}
}

func (st *Stantion) Run(port int) {
	http.HandleFunc("/stantion-ws", st.wsHandler)
	//TODO server there html which i will use as a walkie talkie client
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
