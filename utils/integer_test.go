package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
