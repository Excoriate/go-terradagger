package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

// GetAllFromHost returns all the environment variables from the host
func GetAllFromHost() map[string]string {
	result := make(map[string]string)

	for _, env := range os.Environ() {
		keyValue := strings.Split(env, "=")
		result[keyValue[0]] = utils.RemoveDoubleQuotes(keyValue[1])
	}

	return result
}

// GetAllEnvVarsWithPrefix fetches environment variables that start with the specified prefix
// and returns an error if any of the variables either do not exist or have an empty value.
func GetAllEnvVarsWithPrefix(prefix string) (map[string]string, error) {
	result := make(map[string]string)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]

		if strings.HasPrefix(key, prefix) {
			value := pair[1]
			if value == "" {
				return nil, fmt.Errorf("environment variable %s has an empty value", key)
			}
			result[key] = utils.RemoveDoubleQuotes(value)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no environment variables with the prefix %s found", prefix)
	}

	return result, nil
}

func GetEnvVarByKey(key string, keySensitive bool) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key is empty")
	}

	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s does not exist", key)
	}

	if !keySensitive {
		value = strings.ToLower(value)
	}

	return utils.RemoveDoubleQuotes(value), nil
}
