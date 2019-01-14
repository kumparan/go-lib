package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat32PointerToFloat64(t *testing.T) {
	var f *float32
	assert.Equal(t, float64(0.0), Float32PointerToFloat64(f))
	ff := float32(12.1)
	f = &ff
	assert.Equal(t, float64(ff), Float32PointerToFloat64(f))
	*f = 0
	assert.Equal(t, float64(0), Float32PointerToFloat64(f))
}
