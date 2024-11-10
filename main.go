package main

import (
	peerPackage "concensus/peer"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Arguments required: 'name' and 'port'")
		os.Exit(1)
	}

	name, port := args[0], args[1]
	peer := peerPackage.NewPeer(name, port)

	go peer.StartListening()
	peer.StartApp()
}