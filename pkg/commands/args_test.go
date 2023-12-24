package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatArguments_EmptyCmdArgs(t *testing.T) {
	args := &CmdArgs{}
	formatted := args.FormatArguments()

	assert.Empty(t, formatted, "Formatted arguments should be empty for empty CmdArgs")
}

func TestFormatArguments_FlagType(t *testing.T) {
	args := &CmdArgs{
		{ArgName: "-debug", ArgType: ArgTypeFlag},
	}
	formatted := args.FormatArguments()

	assert.Equal(t, []string{"-debug"}, formatted, "Flag type argument should be formatted correctly")
}

func TestFormatArguments_KeyValueType(t *testing.T) {
	args := &CmdArgs{
		{ArgName: "-var", ArgValue: "foo=bar", ArgType: ArgTypeKeyValue},
	}
	formatted := args.FormatArguments()

	assert.Equal(t, []string{"-var=foo=bar"}, formatted, "Key-Value type argument should be formatted correctly")
}

func TestFormatArguments_ValueType(t *testing.T) {
	args := &CmdArgs{
		{ArgName: "-input", ArgValue: "false", ArgType: ArgTypeValue},
	}
	formatted := args.FormatArguments()

	assert.Equal(t, []string{"-input false"}, formatted, "Value type argument should be formatted correctly")
}

func TestAddNew(t *testing.T) {
	args := &CmdArgs{}
	newArg := CommandArgument{ArgName: "-auto-approve", ArgType: ArgTypeFlag}
	args.AddNew(newArg)

	assert.Equal(t, 1, len(*args), "CmdArgs should contain one element after adding a new argument")
	assert.Equal(t, newArg, (*args)[0], "The added argument should match the newArg")
}
