package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64PointerToFloat64(t *testing.T) {
	var f *float64
	assert.Equal(t, float64(0.0), Float64PointerToFloat64(f))
	ff := float64(12.1)
	f = &ff
	assert.Equal(t, float64(ff), Float64PointerToFloat64(f))
	*f = 0
	assert.Equal(t, float64(0), Float64PointerToFloat64(f))
}
