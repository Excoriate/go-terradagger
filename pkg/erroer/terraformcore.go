package erroer

const ErrTerraformCoreInvalidConfiguration = "Invalid configuration error: "

type ErrTerraformCoreInvalidConfigurationError struct {
	ErrWrapped error
	Msg        string
}

func (e *ErrTerraformCoreInvalidConfigurationError) Error() string {
	if e.ErrWrapped != nil {
		return ErrTerraformCoreInvalidConfiguration + e.Msg + ": " + e.ErrWrapped.Error()
	}

	return ErrTerraformCoreInvalidConfiguration + e.Msg
}

func NewErrTerraformCoreInvalidConfigurationError(errMsg string, err error) *ErrTerraformCoreInvalidConfigurationError {
	return &ErrTerraformCoreInvalidConfigurationError{
		ErrWrapped: err,
		Msg:        errMsg,
	}
}
