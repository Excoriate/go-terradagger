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

	tfLifeCycleCmd := TfLifecycleCMD{}
	tfContainerCfg := &TerraformContainerConfigOptions{
		tfOptions: tfOpts,
		iacConfig: i.Config,
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

	// Native lifecycle command (terraform plan, apply, etc.)
	tfCMDStr, tfCMDStrErr := tfLifeCycleCmd.GetTerraformLifecycleCMDString(&GetTerraformLifecycleCMDStringOptions{
		iacConfig:        i.Config,
		lifecycleCommand: tfLifeCycleCmd.GetPlanCommand(),
		args:             args,
	})

	if tfCMDStrErr != nil {
		return nil, nil, tfCMDStrErr
	}

	tfInitCMDStr, tfCMDInitErr := tfLifeCycleCmd.GenerateTFInitCommandStr(&GenerateTFInitCMDStrOptions{
		iacConfig: i.Config,
		initArgs:  []string{},
	})

	if tfCMDInitErr != nil {
		return nil, nil, tfCMDInitErr
	}

	tfCMDStrShell := terradagger.BuildCMDWithSH(tfCMDStr)
	tfCMDInitStrSHell := terradagger.BuildCMDWithSH(tfInitCMDStr)
	tfCmds := []container.Command{tfCMDStrShell}
	tfInitInjected := []container.Command{tfCMDInitStrSHell}

	td.Log.Info(fmt.Sprintf("running %s plan with the following command: %s", i.Config.GetBinary(), tfCMDStr))

	runtime := tfContainerCfg.getContainerRuntime(td, tfContainerCfg.getContainerImageCfg(td))
	tfContainer := runtime.CreateContainer()
	tfContainer = tfContainerCfg.AddEnvVarsToTerraformContainer(td, runtime, tfContainer)

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
