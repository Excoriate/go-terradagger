package terraform

import (
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/errors"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type InitOptions struct {
	// NoColor is a flag to disable colors in terraform output
	NoColor bool
	// BackendConfigFile is the path to the backend config file
	BackendConfigFile string
	// Upgrade is a flag to upgrade the modules and plugins
	Upgrade bool
}

func (o *InitOptions) validateCMDOptions(terraformDir string) error {
	if o.BackendConfigFile != "" {
		backendConfigFilePath := filepath.Join(terraformDir, o.BackendConfigFile)

		if err := utils.FileExist(backendConfigFilePath); err != nil {
			return &errors.ErrTerraformBackendFileIsNotFound{
				BackendFilePath: o.BackendConfigFile,
				ErrWrapped:      err,
			}
		}

		o.BackendConfigFile = backendConfigFilePath
	}

	return nil
}

// Init Configures a 'terraform init' command and runs it.
func Init(td *terradagger.Client, options *Options, initOptions *InitOptions) error {
	if options == nil {
		options = &Options{}
	}

	if initOptions == nil {
		initOptions = &InitOptions{}
	}

	if err := options.validate(); err != nil {
		return &errors.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    "the options passed to the terraform command are invalid",
		}
	}

	if err := initOptions.validateCMDOptions(options.TerraformDir); err != nil {
		return &errors.ErrTerraformInitFailedToStart{
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
		MountDir:        td.MountDir,
		TerraDaggerCMDs: tfCMDDagger,
	}

	tdOptions.EnvVars = resolveEnvVarsByOptions(options)

	c, err := td.Configure(tdOptions)

	if err != nil {
		return &errors.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    "the terraform init command could not be configured",
		}
	}

	// Run the container.
	return td.Run(c)
}
