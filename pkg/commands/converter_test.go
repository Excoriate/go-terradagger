package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertTerraformVarsToCmdArgs_EmptyMap(t *testing.T) {
	vars := make(map[string]interface{})
	cmdArgs := ConvertTerraformVarsToCmdArgs(vars)

	assert.Empty(t, cmdArgs, "Command arguments should be empty for an empty input map")
}

func TestConvertTerraformVarsToCmdArgs_SingleEntry(t *testing.T) {
	vars := map[string]interface{}{"key1": "value1"}
	cmdArgs := ConvertTerraformVarsToCmdArgs(vars)

	assert.Len(t, cmdArgs, 1, "Command arguments should have one entry for a single map entry")
	assert.Equal(t, "var", cmdArgs[0].ArgName, "ArgName should be 'var'")
	assert.Equal(t, "key1=value1", cmdArgs[0].ArgValue, "ArgValue should be correctly formatted")
	assert.Equal(t, ArgTypeKeyValue, cmdArgs[0].ArgType, "ArgType should be ArgTypeKeyValue")
}

func TestConvertTerraformVarsToCmdArgs_MultipleEntries(t *testing.T) {
	vars := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	cmdArgs := ConvertTerraformVarsToCmdArgs(vars)

	assert.Len(t, cmdArgs, 2, "Command arguments should have two entries for two map entries")
	// Check for each specific key-value pair in the output
	for _, arg := range cmdArgs {
		switch arg.ArgValue {
		case "key1=value1":
			assert.Equal(t, "var", arg.ArgName)
			assert.Equal(t, ArgTypeKeyValue, arg.ArgType)
		case "key2=123":
			assert.Equal(t, "var", arg.ArgName)
			assert.Equal(t, ArgTypeKeyValue, arg.ArgType)
		default:
			t.Errorf("Unexpected ArgValue: %s", arg.ArgValue)
		}
	}
}
