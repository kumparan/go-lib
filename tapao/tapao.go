package tapao

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack"
)

// Marshal default using msgpack, ignore error (positivity)
func Marshal(v interface{}) []byte {
	b, _ := msgpack.Marshal(v)
	return b
}

// Unmarshal default using msgpack
func Unmarshal(b []byte, o interface{}) error {
	return msgpack.Unmarshal(b, o)
}

// SmartUnmarshal tries to unmarshal using msgpack, then try using json when failed
func SmartUnmarshal(b []byte, o interface{}) error {
	err := msgpack.Unmarshal(b, o)
	if err != nil {
		return json.Unmarshal(b, &o)
	}

	return nil
}
