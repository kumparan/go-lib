package utils

import "encoding/json"

// JSON2Map :nodoc:
// DEPRECATED never use this no more
func JSON2Map(j []byte) map[string]interface{} {
	c := make(map[string]interface{})
	_ = json.Unmarshal(j, &c)
	return c
}
