package terraformcore

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type DestroyOptions struct {
	// RefreshOnly is a flag to only refresh the state. Equivalent to
	// terraform plan -refresh-only
	RefreshOnly bool
	// TerraformVarFiles is a list of terraform var files to use
	TerraformVarFiles []string
	// Vars is a list of terraform vars to use
	Vars []TFInputVariable
	// AutoApprove is a flag to auto approve the plan
	AutoApprove bool
}

type DestroyArgs interface {
	GetArgRefreshOnly() []string
	GetArgTerraformVarFiles() []string
	GetArgVars() []string
	GetArgAutoApprove() []string
}

func (po *DestroyOptions) GetArgRefreshOnly() []string {
	if po.RefreshOnly {
		return []string{"-refresh-only"}
	}
	return []string{}
}

func (po *DestroyOptions) GetArgTerraformVarFiles() []string {
	var args []string
	for _, file := range po.TerraformVarFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", file))
	}
	return args
}

func (po *DestroyOptions) GetArgVars() []string {
	var args []string
	for _, v := range po.Vars {
		args = append(args, fmt.Sprintf("-var '%s=%s'", v.Name, utils.EscapeValues(v.Value)))
	}
	return args
}

func (po *DestroyOptions) GetArgAutoApprove() []string {
	if po.AutoApprove {
		return []string{"-auto-approve"}
	}
	return []string{}
}

func (i *IasC) Destroy(td *terradagger.TD, tfOpts TfGlobalOptions, options DestroyArgs, extraArgs []string) (*dagger.Container, container.Runtime, error) {
	if err := tfOpts.IsModulePathValid(); err != nil {
		return nil, nil, err
	}

	var args []string
	if options != nil {
		args = utils.MergeSlices(options.GetArgVars(), options.GetArgTerraformVarFiles(), options.GetArgRefreshOnly(), options.GetArgAutoApprove())
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
			Command:     cmdCfg.GetDestroyCommand(),
			CommandArgs: args,
		})
	} else {
		cmdStr = terradagger.BuildTerraformCommand(terradagger.BuildTerraformCommandOptions{
			Binary:      i.Config.GetBinary(),
			Command:     cmdCfg.GetDestroyCommand(),
			CommandArgs: args,
		})
	}

	tfCommandShell := terradagger.BuildCMDWithSH(cmdStr)

	td.Log.Info(fmt.Sprintf("running %s with the following command: %s", i.Config.GetBinary(), cmdStr))

	// Support for custom container image
	var imageCfg container.Image
	if tfOpts.GetCustomContainerImage() != "" {
		td.Log.Warn(fmt.Sprintf("using custom container image: %s", tfOpts.GetCustomContainerImage()))
		imageCfg = container.NewImageConfig(tfOpts.GetCustomContainerImage(), tfOpts.GetTerraformVersion())
	} else {
		imageCfg = container.NewImageConfig(i.Config.GetContainerImage(), tfOpts.GetTerraformVersion())
	}

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

	if tfOpts.IsMirrorAllEnvVarsFromHost() {
		tfContainer = runtime.AddEnvVars(td.Config.GetHostEnvVars(), tfContainer)
	} else {
		if tfOpts.IsAutoDetectAWSKeysFromHost() {
			tfContainer = runtime.AddEnvVars(td.Config.GetAWSEnvVars(), tfContainer)
		}

		if len(tfOpts.GetEnvVarsToInjectByKeyFromHost()) > 0 {
			envVarsToInject := td.Config.GetEnvVarsByKeys(tfOpts.GetEnvVarsToInjectByKeyFromHost())
			tfContainer = runtime.AddEnvVars(envVarsToInject, tfContainer)
		}

		if tfOpts.IsAutoDetectTFVarsFromHost() {
			tfContainer = runtime.AddEnvVars(td.Config.GetTerraformEnvVars(), tfContainer)
		}
	}

	tfCmds := []container.Command{tfCommandShell}
	tfContainer = runtime.AddCommands(tfCmds, tfContainer)

	return tfContainer, runtime, nil
}

func (i *IasC) DestroyE(td *terradagger.TD, tfOpts TfGlobalOptions, options DestroyArgs, extraArgs []string) (string, error) {
	tfInitContainer, runtime, err := i.Destroy(td, tfOpts, options, extraArgs)
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
