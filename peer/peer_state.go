package peer

import (
	ps "concensus/grpc"
	"log"
	"sync"
)

var peerState = NewPeerState() 

type PeerState struct {
	timestamp int64
	state     string
	mu        sync.Mutex
	cond      *sync.Cond
}

const (
	StateWanted  = "Wanted"
	StateHeld    = "Held"
	StateReleased = "Released"
)

func NewPeerState() *PeerState {
	ps := &PeerState{ 
		timestamp: 0, 
		state: StateReleased, 
		cond: sync.NewCond(&sync.Mutex{}),
	}
	return ps
}

func (ps *PeerState) SetWanted() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.state = StateWanted
}

func (ps *PeerState) SetHeld() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.state = StateHeld
}

func (ps *PeerState) SetReleased() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.state = StateReleased
}

func (ps *PeerState) IsHeld() bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.state == StateHeld
}

func (ps *PeerState) IsWanted() bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.state == StateWanted
}

func (peer *Peer) hasHigherPriority(in *ps.AccessRequest) bool {
	peerState.mu.Lock()
	defer peerState.mu.Unlock()
	return in.Timestamp > peerState.timestamp || (in.Timestamp == peerState.timestamp && in.Name > peer.name)
}

func (ps *PeerState) WaitForRelease(in *ps.AccessRequest) {
	peerState.cond.L.Lock()
	defer peerState.cond.L.Unlock()
	log.Printf("Delaying response to %s until I am finished with critical section\n", in.Name)
	peerState.cond.Wait()
}