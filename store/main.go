package main

import (
	"log"

	"github.com/ganners/gossip"
)

func main() {
	// Create a new gossip server
	server, err := gossip.NewServer(
		"store",
		"An in memory, self replicating data store",
		"0.0.0.0", "8001",
		gossip.NewStdoutLogger(),
	)

	if err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
		return
	}

	kvstore := NewInMemoryStore()

	// Add a few endpoints
	server.Handle("store.read", readHandler(kvstore))
	server.Handle("store.write", writeHandler(kvstore))

	// Transmit state when a store subscribes
	server.Handle("node.subscribe.store", stateSendHandler(kvstore))

	// And read it in when a store sends
	server.Handle("store.stateSend", stateReadHandler(kvstore))

	// Run this server until it is signalled to stop
	<-server.Start()

	log.Println("Shutting down")
}
