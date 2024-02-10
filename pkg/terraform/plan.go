package terraform

import (
	"dagger.io/dagger"
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

func Plan(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options PlanOptions) (*dagger.Container, container.Runtime, error) {
	IaacCfg := terraformcore.IacConfigOptions{
		Binary: "terraform",
	}

	tfIaac := terraformcore.IasC{
		Config: &IaacCfg,
	}

	return tfIaac.Plan(td, tfOpts, &terraformcore.PlanOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
	}, []string{})
}

func PlanE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options PlanOptions) (string, error) {
	IaacCfg := terraformcore.IacConfigOptions{
		Binary: "terraform",
	}

	tfIaac := terraformcore.IasC{
		Config: &IaacCfg,
	}

	return tfIaac.PlanE(td, tfOpts, &terraformcore.PlanOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
	}, []string{})
}
