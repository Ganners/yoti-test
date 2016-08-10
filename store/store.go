package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ganners/core_interview/store/pb/store"
	"github.com/gogo/protobuf/proto"
)

const (
	// Defines the initial size of the hashmap, if it tends to grow large then
	// this should be larger
	defaultBucketSize = 1000
)

// KVStore is an interface which can read and write arbitrary slices of bytes
// to keys which are named by slices of bytes
type KVStore interface {
	// Read and Write
	Read(key []byte) ([]byte, error)
	Write(key, value []byte, overwrite bool) error

	// Convert to and from protobuf
	FromProtoBytes(bytes []byte) error
	ToProto() proto.Message
}

var (
	// Errors...
	ErrorNotFound       = errors.New("the requested key could not be found")
	ErrorNotOverwritten = errors.New("Key already exists, not overwriting")

	// Implementation matches interface
	_ KVStore = &InMemoryStore{}
)

// The InMemoryStore is an implementer of the KVStore, it houses all
// of its data inside a hashmap
type InMemoryStore struct {
	sync.RWMutex
	data map[string][]byte
}

// NewInMemoryStore allocates the internal data to a given size and returns a
// store which implements the KVStore
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string][]byte, defaultBucketSize),
	}
}

// Read will look inside the data structure for the key, in the case of an
// error a generic ErrorNotFound is thrown
func (kv *InMemoryStore) Read(key []byte) ([]byte, error) {

	kv.RLock()
	defer kv.RUnlock()

	data, found := kv.data[string(key)]
	if !found {
		return nil, ErrorNotFound
	}
	return data, nil
}

// Write will simply add the value at the given key in the map. If overwrite is
// false and the item already exists, a generic ErrorNotOverwritten is returned
func (kv *InMemoryStore) Write(key []byte, value []byte, overwrite bool) error {

	kv.Lock()
	defer kv.Unlock()

	strKey := string(key)

	if _, exists := kv.data[strKey]; !overwrite && exists {
		return ErrorNotOverwritten
	}

	kv.data[strKey] = value
	return nil
}

// Outside of the Read/Write interface, it's useful for us to be able to
// convert our data into a protobuf so that we can fully replicate the state to
// new store members in the cluster
func (kv *InMemoryStore) ToProto() proto.Message {

	kv.RLock()
	defer kv.RUnlock()

	return &store.StoreState{
		Data: kv.data,
	}
}

// Outside of the Read/Write interface, it's useful for us to be able to
// convert our data into a protobuf so that we can fully replicate the state to
// new store members in the cluster
func (kv *InMemoryStore) FromProtoBytes(bytes []byte) error {
	storeState := &store.StoreState{}
	err := proto.Unmarshal(bytes, storeState)
	if err != nil {
		fmt.Errorf("unable to unmarshal storeState: %s", err)
	}

	kv.Lock()
	defer kv.Unlock()

	// Merge the data from that into this (don't just replace, that might miss
	// something)
	for k, v := range storeState.GetData() {
		kv.data[k] = v
	}
	return nil
}
