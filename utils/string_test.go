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

func TestGetIDFromSlug(t *testing.T) {
	assert.Equal(t, int64(10), GetIDFromSlug("apakah-10"))
	assert.Equal(t, int64(20), GetIDFromSlug("koji-20"))
	assert.Equal(t, int64(0), GetIDFromSlug("keren-20abc"))
}

func TestStringPointer2String(t *testing.T) {
	var s *string
	assert.Equal(t, "", StringPointer2String(s))
	ss := "bengbeng"
	s = &ss
	assert.Equal(t, "bengbeng", StringPointer2String(s))
	*s = ""
	assert.Equal(t, "", StringPointer2String(s))
}

func TestStringPointer2Float64(t *testing.T) {
	var s *string
	assert.Equal(t, float64(0), StringPointer2Float64(s))
	ss := "12.22"
	s = &ss
	assert.Equal(t, float64(12.22), StringPointer2Float64(s))
	*s = ""
	assert.Equal(t, float64(0), StringPointer2Float64(s))
}

func TestStringPointer2Int64(t *testing.T) {
	var s *string
	assert.Equal(t, int64(0), StringPointer2Int64(s))
	ss := "12"
	s = &ss
	assert.Equal(t, int64(12), StringPointer2Int64(s))
	*s = ""
	assert.Equal(t, int64(0), StringPointer2Int64(s))
}

func TestArrayStringPointer2ArrayInt64(t *testing.T) {
	var ps1,ps2 *string
	s1 := "123"
	s2 := "321"
	ps1 = &s1
	ps2 = &s2

	s := &[]*string{ps1,ps2}
	as :=  ArrayStringPointer2ArrayInt64(s)

	assert.Contains(t, as, int64(123))
	assert.Contains(t, as, int64(321))

}
