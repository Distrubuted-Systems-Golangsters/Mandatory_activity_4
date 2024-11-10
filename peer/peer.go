package peer

import (
	"bufio"
	ps "concensus/grpc"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

type Peer struct {
	ps.UnimplementedPeerServiceServer
	peers map[string]*PeerConnection
	name string
	port string
}

func NewPeer(name, port string) *Peer {
	return &Peer{
		name:      name,
		port:      port,
		peers:     make(map[string]*PeerConnection),
	}
}

func (peer *Peer) StartListening() {
	lis, err := net.Listen("tcp", peer.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	ps.RegisterPeerServiceServer(grpcServer, peer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (peer *Peer) StartApp() string {
	for {
		reader := bufio.NewReader(os.Stdin)
		enteredString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console: %v\n", err)
		}
		enteredString = strings.ToLower(strings.Trim(enteredString, "\r\n"))

		if enteredString == "get" {
			log.Printf("Requesting access to Critical Section\n")
			peer.requestAccessToCritical()
		}
	}
}

func (peer *Peer) requestAccessToCritical() {
	peer.registerOtherPeers()

	peerState.SetWanted()

	// Uncomment to be able to simulate 2 peers in wanted state
	// time.Sleep(10 * time.Second)

	var wg sync.WaitGroup
	peerState.timestamp++;

	for _, otherPeer := range peer.peers {
		wg.Add(1)

		go func(otherPeer *PeerConnection) {
			defer wg.Done()

			_, err := otherPeer.client.RequestAccessFromPeers(context.Background(), &ps.AccessRequest{
				Name: peer.name,
				Timestamp: peerState.timestamp,
			})
			
			if err != nil {
				log.Fatalf("%s is not responding { %v }\n", otherPeer.name, err)
			}
		}(otherPeer)
	}

	wg.Wait()

	peerState.SetHeld()
	CriticalSection()
	peerState.SetReleased()

	peerState.cond.Broadcast()
}