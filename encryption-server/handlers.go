package main

import (
	"errors"
	"fmt"

	"github.com/Ganners/yoti-test/encryption-server/pb/encryption"
	"github.com/Ganners/yoti-test/store/pb/store"
	"github.com/ganners/gossip"
	"github.com/ganners/gossip/pb/envelope"
	"github.com/gogo/protobuf/proto"
)

// Read handler handles encrypted reads
func readHandler(server *gossip.Server, request envelope.Envelope) error {

	server.Logger.Debugf("Received read request")

	// Unmarshal the message
	readRequest := &encryption.EncryptedReadRequest{}
	err := proto.Unmarshal(request.EncodedMessage, readRequest)
	if err != nil {
		return fmt.Errorf("unable to unmarshal read request: %+v", err)
	}

	id := readRequest.Id
	key := readRequest.Key

	if len(id) == 0 || len(key) == 0 {
		return errors.New("id or key length is 0")
	}

	server.Logger.Debugf("Broadcasting to read from store")
	err, rsp, timeout := server.BroadcastAndWaitForResponse(
		"store.read",
		&store.ReadRequest{
			Key: id,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to broadcast request: %+v", err)
	}

	select {
	case <-timeout:
		return errors.New("Timed out waiting for response")
	case responseEnvelope := <-rsp:
		// Unmarshal response
		server.Logger.Debugf("Received from store, decrypting")
		storeReadRsp := &store.ReadResponse{}
		err := proto.Unmarshal(responseEnvelope.EncodedMessage, storeReadRsp)
		if err != nil {
			return fmt.Errorf("unable to unmarshal read request: %+v", err)
		}

		// Decrypt
		decrypted, err := decrypt(storeReadRsp.Value, key)
		if err != nil {
			return fmt.Errorf("unable to decrypt response: %s", err)
		}

		// Broadcast response if there's a receipt
		server.Logger.Debugf("Returning read response")
		readResponse := &encryption.EncryptedReadResponse{
			Plaintext: decrypted,
		}
		headers := request.GetHeaders()
		if headers != nil && len(headers.Receipt) > 0 {
			_, err := server.Broadcast(request.Headers.Receipt, readResponse, int32(envelope.Envelope_RESPONSE))
			return err
		}

		return nil
	}
	return nil
}

// Write handler handles encrypted writes
func writeHandler(server *gossip.Server, request envelope.Envelope) error {

	server.Logger.Debugf("Received write request")

	// Unmarshal the message
	writeRequest := &encryption.EncryptedWriteRequest{}
	err := proto.Unmarshal(request.EncodedMessage, writeRequest)
	if err != nil {
		return fmt.Errorf("unable to unmarshal read request: %+v", err)
	}

	id := writeRequest.Id
	plaintext := writeRequest.Plaintext

	if len(id) == 0 || len(plaintext) == 0 {
		return errors.New("id or plaintext length is 0")
	}

	key := randomBytes(64)
	encrypted, err := encrypt(plaintext, key)
	if err != nil {
		return fmt.Errorf("could not encrypt plaintext: %s", err)
	}

	// Async write request and return key
	server.Logger.Debugf("Broadcasting write to store")
	_, err = server.Broadcast(
		"store.write",
		&store.WriteRequest{
			Key:       id,
			Value:     encrypted,
			Overwrite: true,
		},
		int32(envelope.Envelope_ASYNC_REQUEST),
	)
	if err != nil {
		return fmt.Errorf("unable to broadcast request: %+v", err)
	}

	writeResponse := &encryption.EncryptedWriteResponse{
		Key: key,
	}
	server.Logger.Debugf("About to send back response")
	headers := request.GetHeaders()
	if headers != nil && len(headers.Receipt) > 0 {
		server.Logger.Debugf("Broadcasting response...")
		_, err := server.Broadcast(request.Headers.Receipt, writeResponse, int32(envelope.Envelope_RESPONSE))
		return err
	}

	return nil
}
