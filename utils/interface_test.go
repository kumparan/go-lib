package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceToInt64(t *testing.T) {
	someInteger := 8
	result := InterfaceToInt64(interface{}(someInteger))
	assert.EqualValues(t, someInteger, result)
}
