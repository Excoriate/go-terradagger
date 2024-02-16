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

type ErrTerraformCoreInvalidArgumentError struct {
	BaseError
}

const ErrTerraformCoreInvalidArgumentErrorPrefix = "Invalid argument passed to Terraform Core APIs"

func NewErrTerraformCoreInvalidArgumentError(errMsg string, err error) *ErrTerraformCoreInvalidArgumentError {
	return &ErrTerraformCoreInvalidArgumentError{
		BaseError: BaseError{
			ErrWrapped: err,
			ErrMsg:     fmt.Sprintf("%s: %s", ErrTerraformCoreInvalidArgumentErrorPrefix, errMsg),
		},
	}
}
