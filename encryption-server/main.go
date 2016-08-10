package main

import (
	"log"

	"github.com/ganners/gossip"
	"github.com/ganners/gossip/pb/envelope"
)

func main() {
	// Create a new gossip server
	server, err := gossip.NewServer(
		"encryption-server",
		"A store which handles encryption/decryption of data",
		"0.0.0.0", "8002",
		gossip.NewStdoutLogger(),
	)

	if err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
		return
	}

	server.Handle("encryption-server.read", readHandler)
	server.Handle("encryption-server.write", readHandler)

	// Run this server until it is signalled to stop
	<-server.Start()

	log.Println("Shutting down")
}

// @TODO(mark): Implement handlers and encryption/decryption/key
// generation logic

// Read handler handles encrypted reads
func readHandler(server *gossip.Server, request envelope.Envelope) error {
	return nil
}

// Read handler handles encrypted writes
func writeHandler(server *gossip.Server, request envelope.Envelope) error {
	return nil
}
