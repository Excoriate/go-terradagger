package terraform

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/commands"
)

// convertInputVars converts the terraform input vars to cmd args
// that are compatible with the dagger client.
func convertInputVars(tfVars map[string]interface{}) ([]commands.CommandArgument, error) {
	if tfVars == nil {
		return nil, fmt.Errorf("failed to convert terraform input vars to cmd args: terraform vars are nil")
	}

	if len(tfVars) == 0 {
		return nil, fmt.Errorf("failed to convert terraform input vars to cmd args: terraform vars are empty")
	}

	var cmdArgs []commands.CommandArgument
	for key, value := range tfVars {
		cmdArgs = append(cmdArgs, commands.CommandArgument{
			ArgName:  "-var",
			ArgValue: fmt.Sprintf("%s=%v", key, value), // Arguments will be concatenated as "key=value"
			ArgType:  commands.ArgTypeKeyValue,
		})
	}

	return cmdArgs, nil
}
