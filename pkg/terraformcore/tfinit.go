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

	tfLifeCycleCmd := TfLifecycleCMD{}
	tfContainerCfg := &TerraformContainerConfigOptions{
		tfOptions: tfOpts,
		iacConfig: i.Config,
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

	// Native lifecycle command (terraform plan, apply, etc.)
	tfCMDStr, tfCMDStrErr := tfLifeCycleCmd.GetTerraformLifecycleCMDString(&GetTerraformLifecycleCMDStringOptions{
		iacConfig:        i.Config,
		lifecycleCommand: tfLifeCycleCmd.GetInitCommand(),
		args:             args,
	})

	if tfCMDStrErr != nil {
		return nil, nil, tfCMDStrErr
	}

	tfCommandShell := terradagger.BuildCMDWithSH(tfCMDStr)
	td.Log.Info(fmt.Sprintf("running %s plan with the following command: %s", i.Config.GetBinary(), tfCMDStr))

	runtime := tfContainerCfg.getContainerRuntime(td, tfContainerCfg.getContainerImageCfg(td))
	tfContainer := runtime.CreateContainer()
	tfContainer = tfContainerCfg.AddEnvVarsToTerraformContainer(td, runtime, tfContainer)

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
