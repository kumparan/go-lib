package tapao

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

type TestStruct struct {
	Int64   int64   `json:"int64"`
	Float64 float64 `json:"float64"`
	Bool    bool    `json:"bool"`
	String  string  `json:"string"`
}

var testStruct = TestStruct{
	Int64:   222,
	Float64: 333.333,
	Bool:    true,
	String:  "ferdian the best",
}

func TestMarshal(t *testing.T) {
	in := testStruct
	b, _ := msgpack.Marshal(in)
	assert.Equal(t, b, Marshal(in))
}

func TestUnmarshal(t *testing.T) {
	in := testStruct
	b := Marshal(in)

	var out map[string]interface{}
	err := Unmarshal(b, &out)
	assert.NoError(t, err)
	assert.EqualValues(t, in, out)
}

func TestSmartUnmarshal(t *testing.T) {
	in := testStruct
	var out TestStruct

	t.Run("msgpack", func(t *testing.T) {
		bmsgpack := Marshal(in)
		err := SmartUnmarshal(bmsgpack, &out)
		assert.NoError(t, err)
		assert.EqualValues(t, in, out)
	})

	t.Run("json", func(t *testing.T) {
		bjson, _ := json.Marshal(in)
		err := SmartUnmarshal(bjson, &out)
		assert.NoError(t, err)
		assert.EqualValues(t, in, out)
	})
}
