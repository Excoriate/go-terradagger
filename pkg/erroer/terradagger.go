package erroer

import (
	"fmt"
)

type ErrTerraDaggerInvalidArgumentError struct {
	BaseError // Embedding BaseError
}

const ErrTerraDaggerInvalidArgumentErrorPrefix = "Invalid argument errors"

// NewErrTerraDaggerInvalidArgumentError creates a new ErrTerraDaggerInvalidArgumentError. It utilizes the BaseError struct for common functionality.
func NewErrTerraDaggerInvalidArgumentError(errMsg string, err error) *ErrTerraDaggerInvalidArgumentError {
	return &ErrTerraDaggerInvalidArgumentError{
		BaseError: BaseError{
			ErrWrapped: err,
			ErrMsg:     fmt.Sprintf("%s: %s", ErrTerraDaggerInvalidArgumentErrorPrefix, errMsg),
		},
	}
}
