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


	//need to handle ws connection closing and where to close connection
	if err != nil {
		http.Error(w, "Failed to make websocket connection", 500)
	}

	id := uuid.New().String()
	peer := NewRawPeer(id)
	peer.ws = conn

	st.InitPeer(peer)

	
}

func (st *Stantion) InitPeer(peer *OnePeer) {
	go func() {
		st.peerWsMsgHandler(peer)
	}()

}

func (st *Stantion) peerWsMsgHandler(peer *OnePeer) {

	for {
		var msg WsMessage
		err := peer.ws.ReadJSON(&msg)

		if err != nil {
			//TODO handle error there
		}
		switch msg.Type {
		case MSG_SDP:
			//Add peer to stantion

			newPeer, err := peer.OfferHandler(msg.Data)
			//Add there stantion of peer

			if err != nil {
				//TODO return error
				return
			}
			newPeer.Stantion = st
			go func() {
				newPeer.RunStateControl()
			}()

			err =st.Add(newPeer)
			if err!=nil{
				//TODO return via websocket
				fmt.Println(err)
			}

		case MSG_ICE:
			peer.IceHandler(msg.Ice)
		case MSG_ERROR:
			fmt.Println(msg.Data)
			return
		case MSG_SPEAK_END:
			peer.EndSpeak()
		case MSG_SPEAK_START:
			peer.StartSpeak()

		}

	}
}

func (st *Stantion) Run(port int) {
	http.HandleFunc("/stantion-ws", st.wsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client.html")
	})

	//TODO server there html which i will use as a walkie talkie client
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
