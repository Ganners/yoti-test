package main

import (
	"log"

	"github.com/ganners/gossip"
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
	server.Handle("encryption-server.write", writeHandler)

	// Run this server until it is signalled to stop
	<-server.Start()

	log.Println("Shutting down")
}
