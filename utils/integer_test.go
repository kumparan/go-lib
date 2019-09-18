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

func TestIsSameSliceIgnoreOrder(t *testing.T) {
	a := []int64{2, 1, 3}
	b := []int64{2, 1, 3}
	assert.True(t, IsSameSliceIgnoreOrder(a, b))
	a = []int64{2, 1, 3}
	b = []int64{1, 2, 3}
	assert.True(t, IsSameSliceIgnoreOrder(a, b))
	a = []int64{2, 1, 3, 4}
	b = []int64{1, 2, 3}
	assert.False(t, IsSameSliceIgnoreOrder(a, b))
	a = []int64{}
	b = []int64{}
	assert.True(t, IsSameSliceIgnoreOrder(a, b))
}

func TestLimitInt64(t *testing.T) {
	a := int64(5)
	b := int64(10)
	c := int64(15)
	d := int64(-1)

	assert.Equal(t, a, LimitInt64(b, a))
	assert.Equal(t, b, LimitInt64(b, c))
	assert.Equal(t, a, LimitInt64(d, a))
	assert.NotEqual(t, c, LimitInt64(c, a))

}
