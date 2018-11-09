package token

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomBytes(t *testing.T) {
	randomBytes, _ := GenerateRandomBytes(32)

	// The length of bytes has the same value with the input
	assert.Equal(t, 32, len(randomBytes))

	// TODO: Test the consistency of random generated distribution
}

func TestGenerateRandomString(t *testing.T) {
	randomString, _ := GenerateRandomString(64)

	// The length of string has the same value with the input
	assert.Equal(t, 64, len(randomString))

	// TODO: Test the consistency of random generated distribution
}

func TestGenerateRandomStringURLSafe(t *testing.T) {
	randomString, _ := GenerateRandomStringURLSafe(64)

	// Check wether the Toke URL is safe or not
	_, err := base64.URLEncoding.DecodeString(randomString)
	if err != nil {
		t.Fatal("Token not URL Safe")
	}

	// TODO: Test the consistency of random generated distribution
}
