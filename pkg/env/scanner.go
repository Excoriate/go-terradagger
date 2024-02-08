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

	if prefix == "" {
		return nil, fmt.Errorf("prefix cannot be empty")
	}

	allEnvs := GetAllFromHost()

	for key, value := range allEnvs {
		if strings.HasPrefix(key, prefix) {
			result[key] = value
		}
	}

	return result, nil
}

func GetEnvVarByKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	allEnvs := GetAllFromHost()

	if val, ok := allEnvs[key]; ok {
		return val, nil
	}

	return "", fmt.Errorf("environment variable %s does not exist", key)
}
