package stantion

import (
	"sync"
)

// Stantion model half duplex only one peer can speak at same time
// TODO speaking functions not implmented now
type Stantion struct {
	Peers      map[string]*OnePeer
	mu         sync.RWMutex
	speaking   string
	speakingMu sync.Mutex
}

func NewStantion() *Stantion {
	return &Stantion{
		Peers: make(map[string]*OnePeer),
	}
}

func (st *Stantion) Add(peer *OnePeer) error {
	st.mu.Lock()
	st.Peers[peer.Id] = peer
	st.mu.Unlock()
	return nil
}

func (st *Stantion) Get() {

}

func (st *Stantion) Remove() {

}

func (st *Stantion) All() []*OnePeer {
	st.mu.RLock()
	defer st.mu.RUnlock()

	peers := make([]*OnePeer, 0, len(st.Peers))
	for _, p := range st.Peers {
		peers = append(peers, p)
	}
	return peers
}
func (st *Stantion) PeerCount() int {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return len(st.Peers)
}

func (st *Stantion) SpeakStart(speakerId string) {
	st.speakingMu.Lock()
	if st.speaking != "" {
		st.speakingMu.Unlock()
		return
	}
	st.speaking = speakerId
	st.speakingMu.Unlock()
}

func (s *Stantion) CanSpeak(speakerId string) bool {
	s.speakingMu.Lock()
	defer s.speakingMu.Unlock()
	return s.speaking == speakerId || s.speaking == ""
}
