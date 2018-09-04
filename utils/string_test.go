package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardizeSpaces(t *testing.T) {
	s := `test
	1
	2
	3`

	assert.Equal(t, "test 1 2 3", StandardizeSpaces(s))
}

func TestGenerateSlug(t *testing.T) {
	title1 := "are you okay 100%"
	title2 := "bracket {this}"

	assert.Equal(t, "are-you-okay-100", GenerateSlug(title1))
	assert.Equal(t, "bracket-this", GenerateSlug(title2))
}

func TestString2Bool(t *testing.T) {
	assert.Equal(t, true, String2Bool("true"))
	assert.Equal(t, false, String2Bool("false"))
	assert.Equal(t, false, String2Bool("0"))
	assert.Equal(t, true, String2Bool("1"))
	assert.Equal(t, false, String2Bool("bebek"))
}

func TestString2Int64(t *testing.T) {
	assert.Equal(t, int64(10), String2Int64("10"))
	assert.Equal(t, int64(20), String2Int64("20"))
	assert.Equal(t, int64(0), String2Int64("20abc"))
}

func TestString2Int64WithDefault(t *testing.T) {
	assert.Equal(t, int64(10), String2Int64WithDefault("10", 0))
	assert.Equal(t, int64(20), String2Int64WithDefault("20", 0))
	assert.Equal(t, int64(999), String2Int64WithDefault("20abc", 999))
}
