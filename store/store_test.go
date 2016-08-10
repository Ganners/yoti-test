package main

import (
	"reflect"
	"testing"
)

// Fairly trivial set of tests operating on a store. Will at least
// give us race detection at test time 👍
func TestInMemoryStoreReadWrite(t *testing.T) {

	store := NewInMemoryStore()

	for _, test := range []struct {
		inputKey           []byte
		inputValue         []byte
		overwrite          bool
		expectedOutput     []byte
		expectedWriteError error
		expectedReadError  error
		skipWrite          bool
	}{
		{
			inputKey:           []byte("key 1"),
			inputValue:         []byte("value 1"),
			expectedOutput:     []byte("value 1"),
			overwrite:          false,
			expectedWriteError: nil,
			expectedReadError:  nil,
		},
		{
			inputKey:           []byte("key 1"), // Duplicated
			inputValue:         []byte("value 2"),
			expectedOutput:     []byte("value 1"),
			overwrite:          false,
			expectedWriteError: ErrorNotOverwritten,
			expectedReadError:  nil,
		},
		{
			inputKey:           []byte("key 1"), // Duplicated
			inputValue:         []byte("value 2"),
			expectedOutput:     []byte("value 2"),
			overwrite:          true, // Will let us overwrite
			expectedWriteError: nil,
			expectedReadError:  nil,
		},
		{
			inputKey:           []byte("key 2"),
			inputValue:         []byte("value 4"),
			expectedOutput:     []byte("value 4"),
			overwrite:          true,
			expectedWriteError: nil,
			expectedReadError:  nil,
		},
		{
			inputKey:           []byte("key 3"),
			inputValue:         []byte("value 5"),
			expectedOutput:     []byte("value 5"),
			overwrite:          false,
			expectedWriteError: nil,
			expectedReadError:  ErrorNotFound,

			skipWrite: true, // We won't actually write it
		},
	} {

		// Test the write is as we expect WRT the error
		if !test.skipWrite {
			err := store.Write(test.inputKey, test.inputValue, test.overwrite)
			if err != test.expectedWriteError {
				t.Errorf("Expected write error of %v, got %v", test.expectedWriteError, err)
			}
		}

		// Test the read is as we expect WRT the error
		data, err := store.Read(test.inputKey)
		if err != test.expectedReadError {
			t.Errorf("Expected read error of '%v', got '%v'", test.expectedReadError, err)
		}

		// If we aren't expecting an error, check what we wrote in
		// matches what we want out
		if test.expectedReadError == nil {
			if !reflect.DeepEqual(data, test.expectedOutput) {
				t.Errorf("Expected read output of '%s', got '%s'", string(test.expectedOutput), string(data))
			}
		}
	}
}
