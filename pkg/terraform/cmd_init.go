package terraform

import (
	"fmt"
	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/errors"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/utils"
	"path/filepath"
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

	tfInitCMD := commands.GetTerraformCommand("init", nil)
	tfInitCMD.OmitBinaryNameInCommand = true

	if initOptions.NoColor {
		td.Logger.Info("The option no-color is set to true")
		tfInitCMD, _ = commands.AddArgsToCommand(tfInitCMD, []commands.Args{
			{
				ArgName:  "no-color",
				ArgValue: "",
			},
		})
	}

	if initOptions.BackendConfigFile != "" {
		td.Logger.Info(fmt.Sprintf("The option backend-config is set to %s", initOptions.BackendConfigFile))
		tfInitCMD, _ = commands.AddArgsToCommand(tfInitCMD, []commands.Args{
			{
				ArgName:  "backend-config",
				ArgValue: initOptions.BackendConfigFile,
			},
		})
	}

	if initOptions.Upgrade {
		td.Logger.Info("The option upgrade is set to true")
		tfInitCMD, _ = commands.AddArgsToCommand(tfInitCMD, []commands.Args{
			{
				ArgName:  "upgrade",
				ArgValue: "",
			},
		})
	}

	// tfCMDDagger := commands.ConvertCommandToDaggerFormat(tfInitCMD)
	cmds := []commands.Command{
		tfInitCMD,
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
