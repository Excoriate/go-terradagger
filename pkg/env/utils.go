package env

import "github.com/Excoriate/go-terradagger/pkg/utils"

type Vars map[string]string

// MergeEnvVars merges the environment variables from the specified maps.
func MergeEnvVars(envVars ...Vars) Vars {
  result := make(Vars)

  for _, env := range envVars {
    for key, value := range env {
      if key != "" && value != "" {
        result[key] = utils.RemoveDoubleQuotes(value)
      }
    }
  }

  return result
}
