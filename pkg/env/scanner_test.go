package env

import (
	"os"
	"testing"

	"github.com/Excoriate/go-terradagger/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFromHost(t *testing.T) {
	// Setting an environment variable for test
	testEnvKey := "TEST_ENV_GETALL"
	testEnvValue := "testValue"
	os.Setenv(testEnvKey, testEnvValue)
	defer os.Unsetenv(testEnvKey)

	result := GetAllFromHost()

	assert.NotEmpty(t, result, "Result should not be empty")
	assert.Equal(t, testEnvValue, result[testEnvKey], "The value of the test environment variable should match")
}

func TestGetAllEnvVarsWithPrefix_ExistingPrefix(t *testing.T) {
	// Setting environment variables for test
	prefix := "TEST_ENV_PREFIX_"
	os.Setenv(prefix+"ONE", "First")
	os.Setenv(prefix+"TWO", "Second")
	defer func() {
		os.Unsetenv(prefix + "ONE")
		os.Unsetenv(prefix + "TWO")
	}()

	result, err := GetAllEnvVarsWithPrefix(prefix)

	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, 2, len(result), "There should be two environment variables with the specified prefix")
	assert.Equal(t, "First", utils.RemoveDoubleQuotes(result[prefix+"ONE"]), "The value of the first variable should match")
	assert.Equal(t, "Second", utils.RemoveDoubleQuotes(result[prefix+"TWO"]), "The value of the second variable should match")
}

func TestGetAllEnvVarsWithPrefix_NonExistingPrefix(t *testing.T) {
	result, err := GetAllEnvVarsWithPrefix("NON_EXISTING_PREFIX_")

	assert.Error(t, err, "Should return an error for non-existing prefix")
	assert.Nil(t, result, "Result should be nil for non-existing prefix")
}

func TestGetAllEnvVarsWithPrefix_EmptyValue(t *testing.T) {
	// Setting an environment variable with empty value for test
	testEnvKey := "TEST_ENV_EMPTY"
	os.Setenv(testEnvKey, "")
	defer os.Unsetenv(testEnvKey)

	_, err := GetAllEnvVarsWithPrefix("TEST_ENV_")

	assert.Error(t, err, "Should return an error for environment variable with empty value")
}
