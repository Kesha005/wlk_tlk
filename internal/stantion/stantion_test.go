package stantion

import "testing"

func TestSpeakerStart(t *testing.T) {

	stantion := NewStantion()
	stantion.speaking = ""
	peer := NewRawPeer("id")
	stantion.Add(peer)
	peer.Stantion = stantion
	peer.Speak()

	if !peer.speaking {
		t.Error("Peer must be speaking")
		return
	}
	if stantion.speaking != "id" {
		t.Error("Stantion speaker id not changed")
		return
	}

}

func TestSpeakerDublicate(t *testing.T) {
	stantion := NewStantion()
	stantion.speaking = ""
	peer := NewRawPeer("id")
	stantion.Add(peer)
	peer.Stantion = stantion
	peer.Speak()

	peer2 := NewRawPeer("id2")
	stantion.Add(peer2)
	peer2.Stantion = stantion

	if stantion.PeerCount() != 2 {
		t.Error("Failed to add second peer")
		return
	}

	ok := stantion.CanSpeak(peer2.Id)
	if ok {
		peer2.Speak()
	}

	if peer2.IsSpeaking() {
		t.Error("This must not speak while other speaking")
		return
	}

	if stantion.speaking == "id2" {
		t.Error("Peer 2 must not speaking now")
		return
	}

}
