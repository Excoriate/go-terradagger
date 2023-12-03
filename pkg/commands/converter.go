package commands

import "fmt"

// ConvertTerraformVarsToCmdArgs takes a map of terraform vars and converts them to a slice of CommandArgument
func ConvertTerraformVarsToCmdArgs(vars map[string]interface{}) []CommandArgument {
	var cmdArgs []CommandArgument
	for key, value := range vars {
		// Convert each key-value pair to a CommandArgument with type ArgTypeKeyValue.
		cmdArgs = append(cmdArgs, CommandArgument{
			ArgName:  "var",
			ArgValue: fmt.Sprintf("%s=%v", key, value),
			ArgType:  ArgTypeKeyValue,
		})
	}
	return cmdArgs
}
