package terraformcore

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/container"

	"dagger.io/dagger"

	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

func (i *IasC) Init(td *terradagger.TD, tfOpts TfGlobalOptions, options InitArgs, _ []string) (*dagger.Container, container.Runtime, error) {
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

	lfCMD := TfLifecycleCMD{}
	// Native lifecycle command (terraform plan, apply, etc.)
	tfCMDStr, tfCMDStrErr := lfCMD.GetTerraformLifecycleCMDString(&GetTerraformLifecycleCMDStringOptions{
		iacConfig:        i.Config,
		lifecycleCommand: lfCMD.GetInitCommand(),
		args:             args,
	})

	if tfCMDStrErr != nil {
		return nil, nil, tfCMDStrErr
	}

	tfCommandShell := terradagger.BuildCMDWithSH(tfCMDStr)
	td.Log.Info(fmt.Sprintf("running %s plan with the following command: %s", i.Config.GetBinary(), tfCMDStr))

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
