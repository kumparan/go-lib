package utils

import (
	"encoding/json"
)

// Dump to json using json marshal
func Dump(i interface{}) string {
	return string(ToByte(i))
}

// ToByte :nodoc:
func ToByte(i interface{}) []byte {
	bt, _ := json.Marshal(i)
	return bt
}

// InterfaceBytesToInt64 will transform cached value that get from the redis to Int64
func InterfaceBytesToInt64(i interface{}) int64 {
	if i != nil {
		bt := i.([]byte)

		var integer int64

		err := json.Unmarshal(bt, &integer)
		if err != nil {
			return 0
		}
		return integer
	}

	return 0

}
