package terraform

import (
	"fmt"

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

type initOptions interface {
	GetArgNoColor() []string
	GetArgBackendConfigFile() []string
	GetArgUpgrade() []string
}

func (i *InitOptions) GetArgNoColor() []string {
	arg := []string{"-no-color"}
	if i.NoColor {
		return arg
	}

	return []string{}
}

func (i *InitOptions) GetArgBackendConfigFile() []string {
	arg := []string{"-backend-config", i.BackendConfigFile}
	if i.BackendConfigFile != "" {
		return arg
	}

	return []string{}
}

func (i *InitOptions) GetArgUpgrade() []string {
	arg := []string{"-upgrade"}
	if i.Upgrade {
		return arg
	}

	return []string{}
}

func Init(td *terradagger.TD, tfOpts TfGlobalOptions, options initOptions) (*dagger.Container, container.Runtime, error) {
	if err := tfOpts.IsModulePathValid(); err != nil {
		return nil, nil, err
	}

	var args []string
	if options != nil {
		args = utils.MergeSlices(options.GetArgUpgrade(), options.GetArgNoColor(), options.GetArgBackendConfigFile())
	}

	if err := tfOpts.ModulePathHasTerraformCode(); err != nil {
		return nil, nil, err
	}

	cmdCfg := NewTerraformCommand()

	tfCommandStr := terradagger.BuildCommand(cmdCfg.GetEntryPoint(),
		cmdCfg.GetInitCommand(), args)

	tfCommandShell := terradagger.BuildCMDWithSH(tfCommandStr)

	td.Log.Info(fmt.Sprintf("running terraform init with the following command: %s", tfCommandStr))

	imageCfg := container.NewImageConfig("", tfOpts.GetTerraformVersion())

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

func InitE(td *terradagger.TD, tfOpts TfGlobalOptions, options initOptions) (string, error) {
	tfInitContainer, runtime, err := Init(td, tfOpts, options)
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
