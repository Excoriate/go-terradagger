package terraform

import "github.com/Excoriate/go-terradagger/pkg/env"

func resolveEnvVarsByOptions(options *Options) map[string]string {
	// Initialize an empty map for environment variables
	envVars := make(map[string]string)

	// If options is nil, return the empty map immediately
	if options == nil {
		return envVars
	}

	// If we should auto-inject TF_VAR_ prefixed environment variables
	if options.AutoInjectTFVAREnvVars {
		tfEnvVars, err := env.GetAllEnvVarsWithPrefix("TF_VAR_")
		if err == nil && len(tfEnvVars) > 0 {
			// Merge TF_VAR_ environment variables with existing ones
			// Assuming MergeEnvVars prioritizes the second map's values when keys conflict
			envVars = env.MergeEnvVars(envVars, tfEnvVars)
		}
	}

	return envVars
}
