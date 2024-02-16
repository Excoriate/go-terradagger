package terragrunt

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/terraformcore"
)

type ApplyOptions struct {
	// RefreshOnly is a flag to only refresh the state. Equivalent to
	// terraform plan -refresh-only
	RefreshOnly bool
	// TerraformVarFiles is a list of terraform var files to use
	TerraformVarFiles []string
	// Vars is a list of terraform vars to use
	Vars []terraformcore.TFInputVariable
	// AutoApprove is a flag to auto approve the plan
	AutoApprove bool
}

func Apply(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options ApplyOptions, _ terraformcore.TerragruntConfig) (*dagger.Container, container.Runtime, error) {
	tgRun := terraformcore.NewTerragruntRunner(td, tfOpts, nil)

	return tgRun.RunApply(config.IacToolTerragrunt, &terraformcore.ApplyArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
		AutoApprove:       options.AutoApprove,
		TfGlobalOptions:   tfOpts,
	})
}

func ApplyE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options ApplyOptions, _ terraformcore.TerragruntConfig) (string, error) {
	tgRun := terraformcore.NewTerragruntRunner(td, tfOpts, nil)

	return tgRun.RunApplyE(config.IacToolTerragrunt, &terraformcore.ApplyArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
		AutoApprove:       options.AutoApprove,
		TfGlobalOptions:   tfOpts,
	})
}
