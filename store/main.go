package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/ganners/gossip"
	"github.com/ganners/gossip/pb/envelope"
	"github.com/ganners/yoti-test/store/pb/store"
	"github.com/gogo/protobuf/proto"
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

// ReadHandler takes a store and returns a handler which uses it, it is
// supposed to be called syncronously as a response will be returned
func readHandler(kvstore KVStore) gossip.RequestHandlerFunc {
	return func(server *gossip.Server, request envelope.Envelope) error {

		server.Logger.Debugf("Dealing with read request")

		// Unmarshal the message
		readRequest := &store.ReadRequest{}
		err := proto.Unmarshal(request.EncodedMessage, readRequest)
		if err != nil {
			return fmt.Errorf("unable to unmarshal read request: %+v", err)
		}

		// Look for data
		server.Logger.Debugf("Reading key: %s", readRequest.Key)
		data, err := kvstore.Read(readRequest.Key)
		if err != nil {
			return fmt.Errorf("unable to read key: %+v", err)
		}

		// Generate response payload and broadcast it out
		readResponse := &store.ReadResponse{
			Value: data,
		}
		headers := request.GetHeaders()
		if headers != nil && len(headers.Receipt) > 0 {
			server.Logger.Debugf("Broadcasting response: %s", readRequest.Key)
			_, err := server.Broadcast(request.Headers.Receipt, readResponse, int32(envelope.Envelope_RESPONSE))
			return err
		}

		// It's a little odd if we get this far, it means someone broadcasted
		// an asyncronous read to us...
		return nil
	}
}

// The write handler doesn't return a response, it can be called via async
// broadcasts
func writeHandler(kvstore KVStore) gossip.RequestHandlerFunc {
	return func(server *gossip.Server, request envelope.Envelope) error {

		server.Logger.Debugf("Dealing with write request")

		// Unmarshal the message
		writeRequest := &store.WriteRequest{}
		err := proto.Unmarshal(request.EncodedMessage, writeRequest)
		if err != nil {
			return fmt.Errorf("unable to unmarshal write request: %+v", err)
		}

		if len(writeRequest.Key) == 0 || len(writeRequest.Value) == 0 {
			return errors.New("cannot write empty key or value")
		}

		server.Logger.Debugf("Adding key: %s", writeRequest.Key)
		return kvstore.Write(writeRequest.Key, writeRequest.Value, writeRequest.Overwrite)
	}
}

// When a node registers, we'll send it our state and inevitable update our
// peers
func stateSendHandler(kvstore KVStore) gossip.RequestHandlerFunc {
	return func(server *gossip.Server, request envelope.Envelope) error {
		// Simply broadcast out or state
		_, err := server.Broadcast(
			"store.stateSend",
			kvstore.ToProto(),
			int32(envelope.Envelope_ASYNC_REQUEST),
		)
		return err
	}
}

// When we receive state, just pass it through to the store which can handle
// everything
func stateReadHandler(kvstore KVStore) gossip.RequestHandlerFunc {
	return func(server *gossip.Server, request envelope.Envelope) error {
		return kvstore.FromProtoBytes(request.EncodedMessage)
	}
}
