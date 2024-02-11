package terraformcore

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type TerraformRunner interface {
	RunInit(binary string, options *InitArgsOptions) (*dagger.Container, container.Runtime, error)
	RunInitE(binary string, options *InitArgsOptions) (string, error)
	RunPlan(binary string, options *PlanArgsOptions) (*dagger.Container, container.Runtime, error)
	RunPlanE(binary string, options *PlanArgsOptions) (string, error)
	Apply(binary string, options *ApplyArgsOptions) (*dagger.Container, container.Runtime, error)
	ApplyE(binary string, options *ApplyArgsOptions) (string, error)
	Destroy(binary string, options *DestroyArgsOptions) (*dagger.Container, container.Runtime, error)
	DestroyE(binary string, options *DestroyArgsOptions) (string, error)
}

type TerraformRunnerOptions struct {
	td              *terradagger.TD
	TfGlobalOptions TfGlobalOptions
}

func NewTerraformRunner(td *terradagger.TD, tfGLobalOptions TfGlobalOptions) TerraformRunner {
	return &TerraformRunnerOptions{
		td:              td,
		TfGlobalOptions: tfGLobalOptions,
	}
}

func getIaacConfigByBinary(binary string) *IacConfigOptions {
	return &IacConfigOptions{
		Binary: binary,
	}
}

func (t *TerraformRunnerOptions) RunInit(binary string, args *InitArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Init(t.td, t.TfGlobalOptions, args, []string{})
}

func (t *TerraformRunnerOptions) RunInitE(binary string, args *InitArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.InitE(t.td, t.TfGlobalOptions, args, []string{})
}

func (t *TerraformRunnerOptions) RunPlan(binary string, args *PlanArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Plan(t.td, t.TfGlobalOptions, args, []string{})
}

func (t *TerraformRunnerOptions) RunPlanE(binary string, args *PlanArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.PlanE(t.td, t.TfGlobalOptions, args, []string{})
}

func (t *TerraformRunnerOptions) Apply(binary string, args *ApplyArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Apply(t.td, t.TfGlobalOptions, args, []string{})
}

func (t *TerraformRunnerOptions) ApplyE(binary string, args *ApplyArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.ApplyE(t.td, t.TfGlobalOptions, args, []string{})
}

func (t *TerraformRunnerOptions) Destroy(binary string, args *DestroyArgsOptions) (*dagger.Container, container.Runtime, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.Destroy(t.td, t.TfGlobalOptions, args, []string{})
}

func (t *TerraformRunnerOptions) DestroyE(binary string, args *DestroyArgsOptions) (string, error) {
	tfIaac := IasC{
		Config: getIaacConfigByBinary(binary),
	}

	return tfIaac.DestroyE(t.td, t.TfGlobalOptions, args, []string{})
}
