package commands //nolint:typecheck

// NewTerraDaggerCMD creates a new Command struct for a Terraform command.
func NewTerraDaggerCMD(binary, command string, args TerraDaggerArgs) *TerraDaggerCMD {
	return &TerraDaggerCMD{
		Binary:  binary,
		Command: command,
		Args:    args,
	}
}

// ConvertCommandsToDaggerFormat converts a Command struct to a Dagger-compatible format.
func ConvertCommandsToDaggerFormat(cmds []TerraDaggerCMD) DaggerEngineCMDs {
	var daggerCmds DaggerEngineCMDs
	for _, cmd := range cmds {
		var commandParts []string
		if !cmd.OmitBinaryNameInCommand {
			commandParts = append(commandParts, cmd.Binary)
		}
		commandParts = append(commandParts, cmd.Command)
		commandParts = append(commandParts, cmd.Args...)
		daggerCmds = append(daggerCmds, [][]string{commandParts})
	}
	return daggerCmds
}
