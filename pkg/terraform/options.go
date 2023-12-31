package terraform

import (
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/terradagger"

	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type Options struct {
	// TerraformSRC is the root directory of the terraform code
	TerraformSRC string
	// TerraformModulePath is the directory of the terraform code
	TerraformModulePath string
	// TerraformVersion is the version of terraform to use
	// By default, it'll use the latest version
	TerraformVersion string
	// AutoInjectTFVAREnvVars is a flag to scan the environment variables and inject them into the terraform code
	// The variables that'll be injected are the ones that start with TF_VAR_
	AutoInjectTFVAREnvVars bool
	// AutoInjectEnvVarsFromHost is a flag to scan the environment variables and inject them into the terraform code
	AutoInjectEnvVarsFromHost bool
	// TerraformCustomContainerImage is the custom image to use for the terraform container
	// If it's passed, it'll override the default image hashicorp/terraform
	TerraformCustomContainerImage string
}

type OptionsValidator interface {
	validate() error
}

type CommandOptionsValidator interface {
	validateCMDOptions(options *Options) error
}

func (o *Options) validate() error {
	dirUtils := utils.DirUtils{}
	if o.TerraformModulePath == "" {
		return &erroer.ErrTerraformOptionsAreInvalid{
			Details: "the terraform directory is required, but it was not passed",
		}
	}

	if o.TerraformSRC == "" {
		return &erroer.ErrTerraformOptionsAreInvalid{
			Details: "the terradagger root directory is required, but it was not passed",
		}
	}

	terraformDir := filepath.Join(o.TerraformSRC, o.TerraformModulePath)

	if err := dirUtils.IsValidDir(terraformDir); err != nil {
		return &erroer.ErrTerraformOptionsAreInvalid{
			Details:    "the terraform directory is invalid",
			ErrWrapped: err,
		}
	}

	if o.AutoInjectEnvVarsFromHost && o.AutoInjectTFVAREnvVars {
		return &erroer.ErrTerraformOptionsAreInvalid{
			Details: "the terraform options are invalid. You cannot use both AutoInjectEnvVarsFromHost and AutoInjectTFVAREnvVars at the same time",
		}
	}

	return nil
}

func setDefaultOptions(td *terradagger.TD, options *Options) {
	if options == nil {
		options = &Options{}
	}

	options.TerraformSRC = td.Config.TerraDagger.Paths.Workspace.SRC
	options.TerraformVersion = "1.6.1" // FIXME: Make these defaults configurable.
}
