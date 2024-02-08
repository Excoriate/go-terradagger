package terraformcore

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/container"

	"dagger.io/dagger"

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

type InitArgs interface {
	GetArgNoColor() []string
	GetArgBackendConfigFile() []string
	GetArgUpgrade() []string
}

func (ti *InitOptions) GetArgNoColor() []string {
	arg := []string{"-no-color"}
	if ti.NoColor {
		return arg
	}

	return []string{}
}

func (ti *InitOptions) GetArgBackendConfigFile() []string {
	arg := []string{"-backend-config", ti.BackendConfigFile}
	if ti.BackendConfigFile != "" {
		return arg
	}

	return []string{}
}

func (ti *InitOptions) GetArgUpgrade() []string {
	arg := []string{"-upgrade"}
	if ti.Upgrade {
		return arg
	}

	return []string{}
}

func (i *IasC) Init(td *terradagger.TD, tfOpts TfGlobalOptions, options InitArgs, extraArgs []string) (*dagger.Container, container.Runtime, error) {
	if err := tfOpts.IsModulePathValid(); err != nil {
		return nil, nil, err
	}

	var args []string
	if options != nil {
		args = utils.MergeSlices(options.GetArgUpgrade(), options.GetArgNoColor(), options.GetArgBackendConfigFile())
	}

	if i.Config.GetBinary() == config.IacToolTerraform {
		if err := tfOpts.ModulePathHasTerraformCode(); err != nil {
			return nil, nil, err
		}
	}

	if i.Config.GetBinary() == config.IacToolTerragrunt {
		if err := tfOpts.ModulePathHasTerragruntHCL(); err != nil {
			return nil, nil, err
		}
	}

	cmdCfg := NewTerraformCommandConfig()

	var cmdStr string
	if i.Config.GetBinary() == config.IacToolTerragrunt {
		cmdStr = terradagger.BuildTerragruntCommand(terradagger.BuildTerragruntCommandOptions{
			Binary:      i.Config.GetBinary(),
			Command:     cmdCfg.GetInitCommand(),
			CommandArgs: args,
		})
	} else {
		cmdStr = terradagger.BuildTerraformCommand(terradagger.BuildTerraformCommandOptions{
			Binary:      i.Config.GetBinary(),
			Command:     cmdCfg.GetInitCommand(),
			CommandArgs: args,
		})
	}

	tfCommandShell := terradagger.BuildCMDWithSH(cmdStr)

	td.Log.Info(fmt.Sprintf("running terraform init with the following command: %s", cmdStr))

	imageCfg := container.NewImageConfig(i.Config.GetContainerImage(), tfOpts.GetTerraformVersion())

	td.Log.Info(fmt.Sprintf("container image: %s", i.Config.GetContainerImage()))
	td.Log.Info(fmt.Sprintf("using the image %s for the terraform container", imageCfg.GetImageTerraform()))
	td.Log.Info(fmt.Sprintf("using the version %s for the terraform container", imageCfg.GetVersion()))

	containerCfg := container.Config{
		MountPathAbs:         td.Config.GetWorkspaceAbs(),
		Workdir:              tfOpts.GetModulePath(),
		ContainerImage:       imageCfg,
		KeepEntryPoint:       false,                           // This will override the container's entrypoint with the command we want to run.
		AddPrivateGitSupport: tfOpts.GetEnableSSHPrivateGit(), // Add support for private git repos.
	}

	runtime := container.New(&containerCfg, td)
	tfContainer := runtime.CreateContainer()
	tfCmds := []container.Command{tfCommandShell}
	tfContainer = runtime.AddCommands(tfCmds, tfContainer)

	return tfContainer, runtime, nil
}

func (i *IasC) InitE(td *terradagger.TD, tfOpts TfGlobalOptions, options InitArgs, extraArgs []string) (string, error) {
	tfInitContainer, runtime, err := i.Init(td, tfOpts, options, extraArgs)
	if err != nil {
		return "", err
	}

	out, execErr := runtime.RunAndGetStdout(tfInitContainer)
	if execErr != nil {
		return "", err
	}

	td.Log.Info(out)
	return out, nil
}
