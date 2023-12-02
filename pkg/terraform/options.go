package terraform

import (
  "github.com/Excoriate/go-terradagger/pkg/errors"
  "github.com/Excoriate/go-terradagger/pkg/utils"
)

type Options struct {
  TerraformDir string
  // TerraformVersion is the version of terraform to use
  // By default, it'll use the latest version
  TerraformVersion string
  // AutoInjectTFVAREnvVars is a flag to scan the environment variables and inject them into the terraform code
  // The variables that'll be injected are the ones that start with TF_VAR_
  AutoInjectTFVAREnvVars bool
  // UseAllEnvVarsFromHost is a flag to scan the environment variables and inject them into the terraform code
  UseAllEnvVarsFromHost bool
  // TerraformCustomContainerImage is the custom image to use for the terraform container
  // If it's passed, it'll override the default image hashicorp/terraform
  TerraformCustomContainerImage string
}

type OptionsValidator interface {
  validate() error
}

type CommandOptionsValidator interface {
  validateCMDOptions(terraformDir string) error
}

func (o *Options) validate() error {
  if o.TerraformDir == "" {
    return &errors.ErrTerraformOptionsAreInvalid{
      Details: "the terraform directory is required, but it was not passed",
    }
  }

  if err := utils.IsValidDir(o.TerraformDir); err != nil {
    return &errors.ErrTerraformOptionsAreInvalid{
      Details:    "the terraform directory is invalid",
      ErrWrapped: err,
    }
  }

  if o.UseAllEnvVarsFromHost && o.AutoInjectTFVAREnvVars {
    return &errors.ErrTerraformOptionsAreInvalid{
      Details: "the terraform options are invalid. You cannot use both UseAllEnvVarsFromHost and AutoInjectTFVAREnvVars at the same time",
    }
  }

  return nil
}
