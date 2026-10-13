# Walkie Talkie 🎙️

A browser-based push-to-talk (walkie-talkie) app. Multiple peers connect via WebSocket, audio flows peer-to-peer over WebRTC, and only one peer can speak at a time (half-duplex).

## How It Works

```text
Browser A ──WS signaling──┐
                          ├── Go Server ── tracks peers, enforces half-duplex
Browser B ──WS signaling──┘

Browser A ────────WebRTC audio (Opus/RTP)───────► Browser B
```

- **Go server** handles WebSocket signaling (SDP, ICE) and forwards RTP packets between peers.
- **Browser client** captures the mic, sends audio over WebRTC, and plays remote audio.
- **Half-duplex:** only one peer can transmit at a time.

## Requirements

- Go 1.21+
- A modern browser (Chrome, Firefox, Safari, Edge)
- No external audio libraries needed — the browser handles mic/speaker.

## Install

```bash
git clone https://github.com/Kesha005/wlk-tlk.git
cd wlk-tlk/cmd
go mod download
```

## Run

```bash
go run main.go -port=8070
```

Then open the client in two or more browser tabs/windows:

```text
http://localhost:8070/
```

## Usage

1. Open the page — it auto-connects to the WebSocket server.
2. Click **🔊 Tap to hear audio** once (browser autoplay policy requires a user gesture).
3. Hold **S** on the keyboard, or press and hold the **🎙️** button with mouse/touch, to talk.
4. Release to stop transmitting.
5. Other peers hear you; you hear them when they talk.

A short double-beep plays at the start and end of each transmission (military-style PTT cue).

## Flags

| Flag    | Default | Description              |
| ------- | ------- | ------------------------ |
| `-port` | `8070`  | HTTP/WebSocket server port |

## Project Structure

```text
main.go                       # Entry point, flag parsing
internal/stantion/stantion.go # Peer registry, half-duplex state
internal/stantion/peer.go     # OnePeer: PeerConnection, tracks, forwarding
internal/stantion/ws.go       # WebSocket handler, signaling
client.html                   # Browser UI
```

## Signaling Messages

| Type      | Direction       | Purpose                     |
| --------- | --------------- | --------------------------- |
| `sdp`     | client → server | Client's WebRTC offer       |
| `answer`  | server → client | Server's WebRTC answer      |
| `ice`     | both ways       | ICE candidate exchange      |
| `s_start` | client → server | Peer started speaking       |
| `s_end`   | client → server | Peer stopped speaking       |
| `error`   | server → client | Error notification          |

## Half-Duplex Rules

- `CanSpeak(id)` returns `true` if no one is speaking, or if the same peer is already speaking.
- `SpeakStart(id)` sets the current speaker; `SpeakEnd(id)` clears it.
- While a peer is silent, their incoming RTP is dropped and not forwarded.

## Troubleshooting

**No audio?**
- Click the **🔊 Tap to hear audio** button first — browsers block autoplay.
- Check the browser tab isn't muted.
- Open the browser console for logs.

**Can't connect?**
- Make sure the server is running on the same host/port.
- Check that the WebSocket endpoint `/stantion-ws` is reachable.

**Only one person can talk at a time — is that a bug?**
- No, it's by design (half-duplex, like a real walkie-talkie).

**Choppy audio?**
- Try a wired connection.
- The server forwards RTP without a jitter buffer; on lossy networks, audio may stutter.

## Dependencies
 
- [pion/webrtc](https://github.com/pion/webrtc) — pure Go WebRTC implementation
- [gorilla/websocket](https://github.com/gorilla/websocket) — WebSocket signaling
- [google/uuid](https://github.com/google/uuid) — peer IDs
