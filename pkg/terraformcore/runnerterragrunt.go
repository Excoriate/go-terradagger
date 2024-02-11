package terraformcore

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type TerragruntRunnerOptions struct {
	td              *terradagger.TD
	TfGlobalOptions TfGlobalOptions
	TgConfig        TerragruntConfig
}

func NewTerragruntRunner(td *terradagger.TD, tfGLobalOptions TfGlobalOptions, tgConfig TerragruntConfig) TerraformRunner {
	return &TerragruntRunnerOptions{
		td:              td,
		TfGlobalOptions: tfGLobalOptions,
		TgConfig:        tgConfig,
	}
}

func (tg *TerragruntRunnerOptions) RunInit(binary string, args *InitArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Init(tg.td, tg.TfGlobalOptions, args, []string{})
}

func (tg *TerragruntRunnerOptions) RunInitE(binary string, args *InitArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.InitE(tg.td, tg.TfGlobalOptions, args, []string{})
}

func (tg *TerragruntRunnerOptions) RunPlan(binary string, args *PlanArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Plan(tg.td, tg.TfGlobalOptions, args, []string{})
}

func (tg *TerragruntRunnerOptions) RunPlanE(binary string, args *PlanArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.PlanE(tg.td, tg.TfGlobalOptions, args, []string{})
}

func (tg *TerragruntRunnerOptions) RunApply(binary string, args *ApplyArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Apply(tg.td, tg.TfGlobalOptions, args, []string{})
}

func (tg *TerragruntRunnerOptions) RunApplyE(binary string, args *ApplyArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.ApplyE(tg.td, tg.TfGlobalOptions, args, []string{})
}

func (tg *TerragruntRunnerOptions) RunDestroy(binary string, args *DestroyArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Destroy(tg.td, tg.TfGlobalOptions, args, []string{})
}

func (tg *TerragruntRunnerOptions) RunDestroyE(binary string, args *DestroyArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.DestroyE(tg.td, tg.TfGlobalOptions, args, []string{})
}
