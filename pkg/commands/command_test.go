package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTerraDaggerCMD(t *testing.T) {
	binary := "terraform"
	command := "apply"
	args := TerraDaggerArgs{"-var='foo=bar'", "-auto-approve"}

	terraCmd := NewTerraDaggerCMD(binary, command, args)

	assert.Equal(t, binary, terraCmd.Binary, "Binary should match input")
	assert.Equal(t, command, terraCmd.Command, "Command should match input")
	assert.Equal(t, args, terraCmd.Args, "Args should match input")
}

func TestConvertCommandsToDaggerFormat_EmptySlice(t *testing.T) {
	cmds := []TerraDaggerCMD{}
	daggerCmds := ConvertCommandsToDaggerFormat(cmds)

	assert.Empty(t, daggerCmds, "Dagger commands should be empty for an empty input slice")
}

func TestConvertCommandsToDaggerFormat_NonEmptySlice(t *testing.T) {
	cmds := []TerraDaggerCMD{
		{
			Binary:                  "terraform",
			Command:                 "init",
			Args:                    TerraDaggerArgs{"-backend=false"},
			OmitBinaryNameInCommand: false,
		},
		{
			Binary:                  "terraform",
			Command:                 "apply",
			Args:                    TerraDaggerArgs{"-auto-approve"},
			OmitBinaryNameInCommand: true,
		},
	}
	daggerCmds := ConvertCommandsToDaggerFormat(cmds)

	assert.Len(t, daggerCmds, 2, "Dagger commands should have the same number of entries as the input slice")

	// Assertions for the first command
	assert.Equal(t, []string{"terraform", "init", "-backend=false"}, daggerCmds[0][0], "First command should be correctly formatted")
	// Assertions for the second command
	assert.Equal(t, []string{"apply", "-auto-approve"}, daggerCmds[1][0], "Second command should omit binary name and be correctly formatted")
}
