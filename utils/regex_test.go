package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexEmail(t *testing.T) {
	t.Run("should return true for valid email", func(t *testing.T) {
		validEmail := "user@domain.com"
		assert.True(t, RegexEmail().MatchString(validEmail))
	})

	t.Run("should return false for invalid email", func(t *testing.T) {
		inValidEmail := []string{
			"user-domain.com",
			"user-domaincom",
			"@user-domain.com",
		}
		for _, email := range inValidEmail {
			assert.False(t, RegexEmail().MatchString(email))
		}
	})
}
