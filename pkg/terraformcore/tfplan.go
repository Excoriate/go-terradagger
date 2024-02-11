package terraformcore

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

func (i *IasC) Plan(td *terradagger.TD, tfOpts TfGlobalOptions, options PlanArgs, _ []string) (*dagger.Container, container.Runtime, error) {
	if err := tfOpts.IsModulePathValid(); err != nil {
		return nil, nil, err
	}

	var args []string
	if options != nil {
		args = utils.MergeSlices(options.GetArgVars(), options.GetArgTerraformVarFiles(), options.GetArgRefreshOnly())
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
		lifecycleCommand: lfCMD.GetPlanCommand(),
		args:             args,
	})

	if tfCMDStrErr != nil {
		return nil, nil, tfCMDStrErr
	}

	tfInitCMDStr, tfCMDInitErr := lfCMD.GenerateTFInitCommandStr(&GenerateTFInitCMDStrOptions{
		iacConfig: i.Config,
		initArgs:  []string{},
	})

	if tfCMDInitErr != nil {
		return nil, nil, tfCMDInitErr
	}

	tfCMDStrShell := terradagger.BuildCMDWithSH(tfCMDStr)
	tfCMDInitStrSHell := terradagger.BuildCMDWithSH(tfInitCMDStr)

	td.Log.Info(fmt.Sprintf("running %s plan with the following command: %s", i.Config.GetBinary(), tfCMDStr))

	containerCfg := container.Config{
		MountPathAbs:         td.Config.GetWorkspaceAbs(),
		Workdir:              tfOpts.GetModulePath(),
		ContainerImage:       getContainerImageCfg(td, i.Config, tfOpts),
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

	tfCmds := []container.Command{tfCMDStrShell}
	tfInitInjected := []container.Command{tfCMDInitStrSHell}

	tfContainer = runtime.AddCommands(tfInitInjected, tfContainer)
	tfContainer = runtime.AddCommands(tfCmds, tfContainer)

	return tfContainer, runtime, nil
}

func (i *IasC) PlanE(td *terradagger.TD, tfOpts TfGlobalOptions, options PlanArgs, extraArgs []string) (string, error) {
	tfInitContainer, runtime, err := i.Plan(td, tfOpts, options, extraArgs)
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
