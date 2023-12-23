package terraform

import "github.com/Excoriate/go-terradagger/pkg/commands"

func initCMDDefault() *commands.TerraDaggerCMD {
	// Setting the required terraform init args.
	tfInitArgs := &commands.CmdArgs{}
	tfInitArgs.AddNew(commands.CommandArgument{
		ArgName:  "-input",
		ArgValue: "false",
		ArgType:  commands.ArgTypeFlag,
	})

	tfInitArgs.AddNew(commands.CommandArgument{
		ArgName:  "-backend",
		ArgValue: "false",
		ArgType:  commands.ArgTypeFlag,
	})

	tfInitArgs.AddNew(commands.CommandArgument{
		ArgName:  "-upgrade",
		ArgValue: "false",
		ArgType:  commands.ArgTypeFlag,
	})

	tfInitCMD := commands.NewTerraDaggerCMD("terraform", "init", tfInitArgs.FormatArguments())
	tfInitCMD.OmitBinaryNameInCommand = true

	return tfInitCMD
}
