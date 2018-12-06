package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsInt64(t *testing.T) {
	s := []int64{1, 2, 3, 4, 5}

	assert.Equal(t, true, ContainsInt64(s, int64(1)))
	assert.Equal(t, true, ContainsInt64(s, int64(2)))
	assert.Equal(t, true, ContainsInt64(s, int64(3)))
	assert.Equal(t, true, ContainsInt64(s, int64(4)))
	assert.Equal(t, false, ContainsInt64(s, int64(6)))
	assert.Equal(t, false, ContainsInt64(s, int64(10)))
}

func TestContainsString(t *testing.T) {
	s := []string{"kuda", "horse", "flower"}

	assert.Equal(t, true, ContainsString(s, "kuda"))
	assert.Equal(t, true, ContainsString(s, "horse"))
	assert.Equal(t, true, ContainsString(s, "flower"))
	assert.Equal(t, false, ContainsString(s, "house"))
	assert.Equal(t, false, ContainsString(s, "rainbos"))
}

func TestSliceAtoi(t *testing.T) {
	s := []string{"1", "2", "3"}

	i, err := SliceAtoi(s)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, i)

	s2 := []string{"kuda", "horse", "flower"}
	_, err = SliceAtoi(s2)
	assert.Error(t, err)
}

func TestDifferenceString(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"c", "d", "e"}

	assert.Equal(t, []string{"a", "b", "d", "e"}, DifferenceString(s1, s2))
	assert.Equal(t, []string{"d", "e", "a", "b"}, DifferenceString(s2, s1))
}

func TestDifferenceInt64(t *testing.T) {
	s1 := []int64{1, 2, 3}
	s2 := []int64{3, 4, 5}

	assert.Equal(t, []int64{1, 2, 4, 5}, DifferenceInt64(s1, s2))
	assert.Equal(t, []int64{4, 5, 1, 2}, DifferenceInt64(s2, s1))
}

func TestUniqString(t *testing.T) {
	s := []string{"a", "a", "b", "a", "d", "b"}

	assert.Equal(t, []string{"a", "b", "d"}, UniqString(s))
}

func TestUniqInt64(t *testing.T) {
	s := []int64{1, 1, 2, 4, 2}

	assert.Equal(t, []int64{1, 2, 4}, UniqInt64(s))
}

func TestSlicePointerInt32PointerToSliceInt64(t *testing.T) {
	var i *[]*int32
	var j []int64
	assert.Equal(t, j, SlicePointerInt32PointerToSliceInt64(i))
}
