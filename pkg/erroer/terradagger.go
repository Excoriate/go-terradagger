package erroer

const ErrTerraDaggerInvalidArgumentErrorPrefix = "Invalid argument error: "

type ErrTerraDaggerInvalidArgumentError struct {
	ErrWrapped error
	Msg        string
}

func (e *ErrTerraDaggerInvalidArgumentError) Error() string {
	if e.ErrWrapped != nil {
		return ErrTerraDaggerInvalidArgumentErrorPrefix + e.Msg + ": " + e.ErrWrapped.Error()
	}

	return ErrTerraDaggerInvalidArgumentErrorPrefix + e.Msg
}

func NewErrTerraDaggerInvalidArgumentError(errMsg string, err error) *ErrTerraDaggerInvalidArgumentError {
	return &ErrTerraDaggerInvalidArgumentError{
		ErrWrapped: err,
		Msg:        errMsg,
	}
}
