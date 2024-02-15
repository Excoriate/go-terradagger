package erroer

// BaseError provides common fields and methods for custom errors.
type BaseError struct {
	ErrWrapped error
	ErrMsg     string
}

// Error generates the error message. It can be overridden by embedded types.
func (e *BaseError) Error() string {
	msg := "error"
	if e.ErrMsg != "" {
		msg += ": " + e.ErrMsg
	}
	if e.ErrWrapped != nil {
		if e.ErrMsg != "" {
			msg += "; wrapped: " + e.ErrWrapped.Error()
		} else {
			msg += ": " + e.ErrWrapped.Error()
		}
	}
	return msg
}

// Unwrap returns the wrapped error, if any.
func (e *BaseError) Unwrap() error {
	return e.ErrWrapped
}
