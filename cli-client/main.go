package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ganners/gossip"
	"github.com/ganners/yoti-test/client"
	"github.com/ganners/yoti-test/encryption-server/pb/encryption"
	"github.com/gogo/protobuf/proto"
)

func main() {
	// Create a new gossip server
	server, err := gossip.NewServer(
		"cli-client",
		"The client which launches syncronous client requests",
		"0.0.0.0", "8003",
		gossip.NewSilentLogger(),
	)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

	// Launch a new client
	client := NewCLIClient(server)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Read or write? [w/r]: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Unable to read input: ", err)
				continue
			}
			text = strings.TrimSpace(text)

			if text == "w" {
				fmt.Print("Enter payload to write: ")
				payload, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input: ", err)
					continue
				}
				payload = strings.TrimSpace(payload)

				fmt.Print("Enter id to write value to: ")
				id, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input: ", err)
					continue
				}
				id = strings.TrimSpace(id)

				key, err := client.Store([]byte(id), []byte(payload))
				if err != nil {
					log.Println("Could not store payload: ", err)
					continue
				}
				log.Println("[Success] Your key to retrieve in future is: ", string(key))
			} else if text == "r" {
				fmt.Print("Enter your ID: ")
				id, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input: %s", err)
					continue
				}
				id = strings.TrimSpace(id)

				fmt.Print("Enter your encryption key: ")
				key, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input: ", err)
					continue
				}
				key = strings.TrimSpace(key)

				payload, err := client.Retrieve([]byte(id), []byte(key))
				if err != nil {
					log.Println("Could not retrieve payload: ", err)
					continue
				}
				log.Println("[Success] Your original payload was: ", string(payload))
			}
		}
	}()

	<-server.Start()
	log.Println("Terminating client server")
}

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
