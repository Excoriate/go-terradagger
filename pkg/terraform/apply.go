package terraform

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

func Apply(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options ApplyOptions) (*dagger.Container, container.Runtime, error) {
	tfRun := terraformcore.NewTerraformRunner(td, tfOpts)

	return tfRun.RunApply(config.IacToolTerraform, &terraformcore.ApplyArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
		AutoApprove:       options.AutoApprove,
	})
}

func ApplyE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options ApplyOptions) (string, error) {
	tfRun := terraformcore.NewTerraformRunner(td, tfOpts)

	return tfRun.RunApplyE(config.IacToolTerraform, &terraformcore.ApplyArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
		AutoApprove:       options.AutoApprove,
	})
}
