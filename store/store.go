package main

import (
	"errors"
	"sync"
)

const (
	// Defines the initial size of the hashmap, if it tends to grow large then
	// this should be larger
	defaultBucketSize = 1000
)

// KVStore is an interface which can read and write arbitrary slices of bytes
// to keys which are named by slices of bytes
type KVStore interface {
	Read(key []byte) ([]byte, error)
	Write(key, value []byte, overwrite bool) error
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
