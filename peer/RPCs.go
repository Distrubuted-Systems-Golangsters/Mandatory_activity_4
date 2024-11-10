package peer

import (
	ps "concensus/grpc"
	"context"
	"log"
)

func (peer *Peer) RequestAccessFromPeers(ctx context.Context, in *ps.AccessRequest) (*ps.AccessResponse, error) {
	if peerState.IsHeld() || (peerState.IsWanted() && peer.hasHigherPriority(in)) {
		peerState.WaitForRelease(in)
	}

	peerState.timestamp = max(peerState.timestamp, in.Timestamp) + 1

	log.Printf("Sending allowAccess response to %s\n", in.Name)
	
	return &ps.AccessResponse{ AllowAccess: true }, nil
}