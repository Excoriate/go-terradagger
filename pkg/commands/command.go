package commands

import (
  "fmt"
  "github.com/Excoriate/go-terradagger/pkg/utils"
)

type Command struct {
  Binary                  string
  Command                 string
  Args                    []Args
  OmitBinaryNameInCommand bool
}

// Args represents a command-line argument, consisting of a ArgName and a ArgValue.
type Args struct {
  ArgName  string
  ArgValue string
}

func ConvertMapIntoTerraformVarsOption(sourceMap map[string]interface{}) ([]Args, error) {
  var args []Args
  for key, value := range sourceMap {
    args = append(args, Args{
      ArgName:  fmt.Sprintf("-var=%s=%v", key, value),
      ArgValue: "",
    })
  }
  return args, nil
}

type TerraDaggerCMDs [][][]string

// GetCommand returns the command as a slice of strings, ready to be executed.
func (c *Command) GetCommand() []string {
  var commandParts []string
  if !c.OmitBinaryNameInCommand {
    commandParts = append(commandParts, c.Binary)
  }
  commandParts = append(commandParts, c.Command)
  for _, arg := range c.Args {
    commandParts = append(commandParts, arg.ArgName)
    commandParts = append(commandParts, arg.ArgValue)
  }
  return commandParts
}

// GetTerraformCommand creates a new Command struct for a Terraform command.
func GetTerraformCommand(command string, args []Args) Command {
  return Command{
    Binary:  "terraform",
    Command: command,
    Args:    args,
  }
}

// AddArgsToCommand adds a slice of Args to a Command struct.
func AddArgsToCommand(command Command, args []Args) (Command, error) {
  for _, arg := range args {
    if arg.ArgName == "" {
      return command, fmt.Errorf("the name of an argument cannot be empty in command %s", command.Command)
    }

    argValue := ""
    if arg.ArgValue != "" {
      argValue = fmt.Sprintf("=%s", arg.ArgValue)
    }

    command.Args = append(command.Args, Args{
      ArgName:  arg.ArgName,
      ArgValue: argValue,
    })
  }

  return command, nil
}

// ConvertCommandsToDaggerFormat converts a Command struct to a Dagger-compatible format.
func ConvertCommandsToDaggerFormat(cmds []Command) TerraDaggerCMDs {
  var daggerCmds TerraDaggerCMDs
  for _, cmd := range cmds {
    cleanedCommand := utils.CleanSliceFromValuesThatAreEmpty(cmd.GetCommand())
    daggerCmds = append(daggerCmds, [][]string{cleanedCommand})
  }
  return daggerCmds
}
