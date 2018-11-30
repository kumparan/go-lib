package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt32PointerToInt64(t *testing.T) {
	var i *int32
	assert.Equal(t, int64(0), Int32PointerToInt64(i))
	ii := int32(12)
	i = &ii
	assert.Equal(t, int64(ii), Int32PointerToInt64(i))
	*i = 0
	assert.Equal(t, int64(0), Int32PointerToInt64(i))
}

func TestInt32PointerToInt32(t *testing.T) {
	var i *int32
	assert.Equal(t, int32(0), Int32PointerToInt32(i))
	ii := int32(12)
	i = &ii
	assert.Equal(t, int32(ii), Int32PointerToInt32(i))
	*i = 0
	assert.Equal(t, int32(0), Int32PointerToInt32(i))
}
