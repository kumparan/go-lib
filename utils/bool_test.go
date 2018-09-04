package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmailValid(t *testing.T) {
	assert.Equal(t, false, IsEmailValid("bebek123"))
	assert.Equal(t, false, IsEmailValid("bebek123/gmal"))
	assert.Equal(t, true, IsEmailValid("bebek123@gmal.com"))
}

func TestIsNumeric(t *testing.T) {
	assert.Equal(t, false, IsNumeric("bebek123"))
	assert.Equal(t, true, IsNumeric("123"))
}

func TestBool2String(t *testing.T) {
	assert.Equal(t, strconv.FormatBool(false), Bool2String(false))
	assert.Equal(t, strconv.FormatBool(true), Bool2String(true))
}
