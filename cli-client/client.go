package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/ganners/gossip"
	"github.com/ganners/yoti-test/client"
	"github.com/ganners/yoti-test/encryption-server/pb/encryption"
	"github.com/gogo/protobuf/proto"
)

var (
	// Make sure CLIClient implements client.Client
	_ client.Client = &CLIClient{}
)

// CLIClient takes a Gossip server and will broadcast and wait for responses
// back
type CLIClient struct {
	server *gossip.Server
}

// NewCLIClient simply returns a new CLIClient - pass in a gossip server that
// will launch
func NewCLIClient(server *gossip.Server) *CLIClient {
	return &CLIClient{
		server: server,
	}
}

// Implements the Store part of the interface
func (cli *CLIClient) Store(id, payload []byte) (aesKey []byte, err error) {
	// Execute a request and wait for response back
	writeRequest := &encryption.EncryptedWriteRequest{
		Id:        id,
		Plaintext: payload,
	}
	err, resp, timeout := cli.server.BroadcastAndWaitForResponse(
		"encryption-server.write",
		writeRequest)
	if err != nil {
		return nil, fmt.Errorf("Unable to broadcast retrieve: %s", err)
	}

	select {
	case <-timeout:
		return nil, errors.New("timed out waiting for retrieve response")
	case responseEnvelope := <-resp:
		writeResp := &encryption.EncryptedWriteResponse{}
		err := proto.Unmarshal(responseEnvelope.EncodedMessage, writeResp)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal write request: %+v", err)
		}

		// Return the text
		return writeResp.Key, nil
	}

	return nil, nil
}

// Implements the Retrieve part of the interface
func (cli *CLIClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	// Execute a request and wait for response back
	readRequest := &encryption.EncryptedReadRequest{
		Key: aesKey,
		Id:  id,
	}

	log.Println("Sending read")
	err, resp, timeout := cli.server.BroadcastAndWaitForResponse(
		"encryption-server.read",
		readRequest)
	if err != nil {
		return nil, fmt.Errorf("Unable to broadcast retrieve: %s", err)
	}

	select {
	case <-timeout:
		return nil, errors.New("timed out waiting for retrieve response")
	case responseEnvelope := <-resp:
		log.Println("RECEIVED RESPONSE")
		readRsp := &encryption.EncryptedReadResponse{}
		err := proto.Unmarshal(responseEnvelope.EncodedMessage, readRsp)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal read request: %+v", err)
		}

		// Return the text
		return readRsp.Plaintext, nil
	}

	return nil, nil
}
