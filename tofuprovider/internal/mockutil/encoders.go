package mockutil

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

// MsgPack returns a MessagePack encoding of the given value, or panics if it
// cannot be encoded.
//
// This is a helper for populating expected values in mocks.
func MsgPack(v any) []byte {
	ret, err := msgpack.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

// JSON returns a compact JSON encoding of the given value, or panics if it
// cannot be encoded.
//
// This is a helper for populating expected values in mocks.
func JSON(v any) []byte {
	ret, err := json.Marshal(v)
	if err != nil {
		panic(err.Error())
	}
	return ret
}
