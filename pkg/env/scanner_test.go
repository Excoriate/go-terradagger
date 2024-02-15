package env

import (
	"os"
	"testing"
)

// GetAllFromHost returns all the environment variables from the host
func setEnvVars(t *testing.T, vars map[string]string) func() {
	t.Helper()
	originalVars := make(map[string]string, len(vars))
	for k, v := range vars {
		if originalValue, exists := os.LookupEnv(k); exists {
			originalVars[k] = originalValue
		}
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("Error setting up environment variable %s for test: %v", k, err)
		}
	}
	return func() {
		for k := range vars {
			if originalValue, exists := originalVars[k]; exists {
				_ = os.Setenv(k, originalValue)
			} else {
				_ = os.Unsetenv(k)
			}
		}
	}
}

func TestGetAllFromHost(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{"TEST_ENV_VAR": "test"})
	defer cleanup()

	envs := GetAllFromHost()
	if envs["TEST_ENV_VAR"] != "test" {
		t.Errorf("Expected to find TEST_ENV_VAR with value 'test', found '%s'", envs["TEST_ENV_VAR"])
	}
}

func TestGetAllEnvVarsWithPrefix(t *testing.T) {
	cleanup := setEnvVars(t, map[string]string{"PREFIX_TEST_ENV": "value"})
	defer cleanup()

	envs, err := GetAllEnvVarsWithPrefix("PREFIX_")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(envs) == 0 || envs["PREFIX_TEST_ENV"] != "value" {
		t.Errorf("Expected to find PREFIX_TEST_ENV, got %v", envs)
	}

	_, err = GetAllEnvVarsWithPrefix("")
	if err == nil {
		t.Errorf("Expected error when prefix is empty")
	}
}

func TestGetEnvVarByKey(t *testing.T) {
	key := "SOME_UNIQUE_KEY"
	expectedValue := "someValue"
	cleanup := setEnvVars(t, map[string]string{key: expectedValue})
	defer cleanup()

	value, err := GetEnvVarByKey(key)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if value != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, value)
	}

	_, err = GetEnvVarByKey("NON_EXISTENT_KEY")
	if err == nil {
		t.Errorf("Expected error for non-existent key")
	}

	_, err = GetEnvVarByKey("")
	if err == nil {
		t.Errorf("Expected error for empty key")
	}
}
