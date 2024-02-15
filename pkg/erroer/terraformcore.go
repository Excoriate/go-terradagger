package erroer

import "fmt"

type ErrTerraformCoreInvalidConfigurationError struct {
	BaseError
}

const ErrTerraformCoreInvalidConfigurationPrefix = "Invalid configuration error while using Terraform Core APIs "

func NewErrTerraformCoreInvalidConfigurationError(errMsg string, err error) *ErrTerraformCoreInvalidConfigurationError {
	return &ErrTerraformCoreInvalidConfigurationError{
		BaseError: BaseError{
			ErrWrapped: err,
			ErrMsg:     fmt.Sprintf("%s: %s", ErrTerraformCoreInvalidConfigurationPrefix, errMsg),
		},
	}
}
