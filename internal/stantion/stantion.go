package stantion

import (
	"sync"

	
)

//Stantion model half duplex only one peer can speak at same time
//TODO speaking functions not implmented now
type Stantion struct{
	Peers map[string]*OnePeer
	mu sync.RWMutex
	speaking string 
	speakingMu sync.Mutex
}



func NewStantion()*Stantion{
	return &Stantion{
		Peers: make(map[string]*OnePeer),
	}
}

func (st *Stantion)Add(peer *OnePeer)error{
	st.mu.Lock()	
	st.Peers[peer.Id] = peer
	st.mu.Unlock()
	return nil
}

func (st *Stantion)Get(){
	
}


func (st *Stantion)Remove(){
	
}

func (st *Stantion)All(){
	
}

