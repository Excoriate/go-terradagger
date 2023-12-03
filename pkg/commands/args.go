package commands

import (
	"fmt"
	"strings"
)

// ArgType is an enum for the different terraform argument types.
type ArgType string

const (
	// ArgTypeFlag represents a boolean flag argument.
	ArgTypeFlag ArgType = "Flag"
	// ArgTypeKeyValue represents a key-value pair argument.
	ArgTypeKeyValue ArgType = "KeyValue"
	// ArgTypeValue represents a simple value argument.
	ArgTypeValue ArgType = "Value"
)

// CommandArgument represents a single command-line argument for a Terraform command.
type CommandArgument struct {
	ArgName  string
	ArgValue interface{} // Using interface{} to allow any type for the argument value.
	ArgType  ArgType
}

// CmdArgs represents a slice of CommandArgument that can be used in Terraform commands.
type CmdArgs []CommandArgument

// FormatArguments formats the command arguments into a slice of strings.
func (args *CmdArgs) FormatArguments() []string {
	var formattedArgs []string
	for _, arg := range *args {
		formattedArg := arg.formatArgument()
		if formattedArg != "" {
			formattedArgs = append(formattedArgs, formattedArg)
		}
	}
	return formattedArgs
}

// AddNew adds a new CommandArgument to the CmdArgs.
func (args *CmdArgs) AddNew(arg CommandArgument) {
	*args = append(*args, arg)
}

// formatArgument helps to format  an individual CommandArgument based on its type.
func (arg CommandArgument) formatArgument() string {
	switch arg.ArgType {
	case ArgTypeFlag:
		// Directly return the flag without any value. Terraform boolean flags do not need an explicit value.
		return fmt.Sprintf("-%s", strings.TrimPrefix(arg.ArgName, "-"))
	case ArgTypeKeyValue:
		// Handle key-value pairs, expect ArgValue to be a string for ArgTypeKeyValue. No need to add explicit "true".
		return fmt.Sprintf("-%s=%v", strings.TrimPrefix(arg.ArgName, "-"), arg.ArgValue)
	case ArgTypeValue:
		// ArgTypeValue is for args which are not boolean flags and don't follow the key-value pair structure.
		// For such cases, we provide the arg value directly after the name separated by a space.
		return fmt.Sprintf("-%s %v", strings.TrimPrefix(arg.ArgName, "-"), arg.ArgValue)
	}
	return ""
}
