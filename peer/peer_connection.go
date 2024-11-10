package peer

import (
	ps "concensus/grpc"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PeerInfo struct {
	name string
	port string
}

var peersList = []PeerInfo{
	{name: "peer1", port: ":5000"},
	{name: "peer2", port: ":5001"},
	{name: "peer3", port: ":5002"},
}

type PeerConnection struct {
	client ps.PeerServiceClient
	name string
}

func (peer *Peer) registerOtherPeer(otherPeerName string, otherPeerPort string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(otherPeerPort, opts...)
	if err != nil {
		log.Fatalf("Connection to %s on port%s failed: %v\n", otherPeerName, otherPeerPort, err)
		return
	}

	peer.peers[otherPeerName] = &PeerConnection{
		client: ps.NewPeerServiceClient(conn),
		name: otherPeerName,
	}
}

func (peer *Peer) registerOtherPeers() {
	for _, otherPeer := range peersList {
		if peer.peers[otherPeer.name] == nil && peer.name != otherPeer.name {
			peer.registerOtherPeer(otherPeer.name, otherPeer.port)
		}
	}
}