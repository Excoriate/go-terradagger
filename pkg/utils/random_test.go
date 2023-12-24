package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString_NoSpecialChars_NoUpperLower(t *testing.T) {
	opt := RandomStringOptions{
		Length:            10,
		AllowSpecialChars: false,
		UpperLowerRandom:  false,
	}

	result := GenerateRandomString(opt)
	assert.Len(t, result, opt.Length, "Generated string should have the specified length")
	assert.Regexp(t, "^[a-zA-Z]+$", result, "Generated string should contain only letters")
}

func TestGenerateRandomString_WithSpecialChars(t *testing.T) {
	opt := RandomStringOptions{
		Length:            10,
		AllowSpecialChars: true,
		UpperLowerRandom:  false,
	}

	result := GenerateRandomString(opt)
	assert.Len(t, result, opt.Length, "Generated string should have the specified length")
	assert.Regexp(t, "^[a-zA-Z!@#$%^&*()]+$", result, "Generated string should contain letters and special characters")
}

func TestGenerateRandomString_WithUpperLowerRandom(t *testing.T) {
	opt := RandomStringOptions{
		Length:            10,
		AllowSpecialChars: false,
		UpperLowerRandom:  true,
	}

	result := GenerateRandomString(opt)
	assert.Len(t, result, opt.Length, "Generated string should have the specified length")
	// Testing randomness of upper/lower cases is impractical
}

func TestGenerateRandomName(t *testing.T) {
	retry := 0
	name := GenerateRandomName(retry)
	assert.NotEmpty(t, name, "Generated name should not be empty")
}
