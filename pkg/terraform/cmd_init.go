package terraform

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type InitOptions struct {
	// NoColor is a flag to disable colors in terraform output
	NoColor bool
	// BackendConfigFile is the path to the backend config file
	BackendConfigFile string
	// Upgrade is a flag to upgrade the modules and plugins
	Upgrade bool
}

func (o *InitOptions) validateCMDOptions(options *Options) error {
	if o.BackendConfigFile == "" {
		return nil
	}

	_, err := getTerraformFileRelativePath(options.TerraformRootDir, options.TerraformDir, o.BackendConfigFile)
	if err != nil {
		return &erroer.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    fmt.Sprintf("the backend config file %s is invalid", o.BackendConfigFile),
		}
	}

	return nil
}

// Init Configures a 'terraform init' command and runs it.
func Init(td *terradagger.Client, options *Options, initOptions *InitOptions) error {
	setDefaultOptions(td, options)

	if initOptions == nil {
		initOptions = &InitOptions{}
	}

	if err := options.validate(); err != nil {
		return &erroer.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    "the options passed to the terraform command are invalid",
		}
	}

	if err := initOptions.validateCMDOptions(options); err != nil {
		return &erroer.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    "the init options passed to the terraform command are invalid",
		}
	}

	td.Logger.Info("All the options are valid, and the terraform init command can be started.")

	// Resolving the arguments
	args := &commands.CmdArgs{}

	if initOptions.NoColor {
		args.AddNew(commands.CommandArgument{
			ArgName:  "no-color",
			ArgValue: "true",
			ArgType:  commands.ArgTypeFlag,
		})
	}

	if initOptions.BackendConfigFile != "" {
		args.AddNew(commands.CommandArgument{
			ArgName:  "backend-config",
			ArgValue: initOptions.BackendConfigFile,
			ArgType:  commands.ArgTypeKeyValue,
		})
	}

	if initOptions.Upgrade {
		args.AddNew(commands.CommandArgument{
			ArgName:  "upgrade",
			ArgValue: "true",
			ArgType:  commands.ArgTypeFlag,
		})
	}

	tfInitCMD := commands.NewTerraDaggerCMD("terraform", "init", args.FormatArguments())
	tfInitCMD.OmitBinaryNameInCommand = true

	cmds := []commands.TerraDaggerCMD{
		*tfInitCMD,
	}

	tfCMDDagger := commands.ConvertCommandsToDaggerFormat(cmds)

	tfImage := resolveTerraformImage(options)
	tfVersion := resolveTerraformVersion(options)

	// Configuring the options.
	tdOptions := &terradagger.ClientConfigOptions{
		Image:           tfImage,
		Version:         tfVersion,
		Workdir:         options.TerraformDir,
		MountDir:        td.Paths.MountDirPath,
		TerraDaggerCMDs: tfCMDDagger,
	}

	tdOptions.EnvVars = resolveEnvVarsByOptions(options)

	c, err := td.Configure(tdOptions)

	if err != nil {
		return &erroer.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    "the terraform init command could not be configured",
		}
	}

	// Run the container.
	return td.Run(c)
}
