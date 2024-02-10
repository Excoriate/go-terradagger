package terraform

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/terraformcore"
)

type DestroyOptions struct {
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

func Destroy(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options DestroyOptions) (*dagger.Container, container.Runtime, error) {
	IaacCfg := terraformcore.IacConfigOptions{
		Binary: config.IacToolTerraform,
	}

	tfIaac := terraformcore.IasC{
		Config: &IaacCfg,
	}

	return tfIaac.Destroy(td, tfOpts, &terraformcore.DestroyArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
		AutoApprove:       options.AutoApprove,
	}, []string{})
}

func DestroyE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options DestroyOptions) (string, error) {
	IaacCfg := terraformcore.IacConfigOptions{
		Binary: config.IacToolTerraform,
	}

	tfIaac := terraformcore.IasC{
		Config: &IaacCfg,
	}

	return tfIaac.DestroyE(td, tfOpts, &terraformcore.DestroyArgsOptions{
		RefreshOnly:       options.RefreshOnly,
		TerraformVarFiles: options.TerraformVarFiles,
		Vars:              options.Vars,
		AutoApprove:       options.AutoApprove,
	}, []string{})
}
