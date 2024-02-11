package terragrunt

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/terraformcore"
)

type PlanOptions struct {
	// RefreshOnly is a flag to only refresh the state. Equivalent to
	// terraform plan -refresh-only
	RefreshOnly bool
	// TerraformVarFiles is a list of terraform var files to use
	TerraformVarFiles []string
	// Vars is a list of terraform vars to use
	Vars []terraformcore.TFInputVariable
}

func Plan(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options PlanOptions, _ terraformcore.TerragruntConfig) (*dagger.Container, container.Runtime, error) {
	tgRun := terraformcore.NewTerragruntRunner(td, tfOpts, nil)

	return tgRun.RunPlan(config.IacToolTerragrunt, &terraformcore.PlanArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
	})
}

func PlanE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options PlanOptions, _ terraformcore.TerragruntConfig) (string, error) {
	tgRun := terraformcore.NewTerragruntRunner(td, tfOpts, nil)

	return tgRun.RunPlanE(config.IacToolTerragrunt, &terraformcore.PlanArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
	})
}
