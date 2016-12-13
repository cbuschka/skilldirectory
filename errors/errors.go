package errors

// Implements error interface
type BadRequestError struct {
	ErrorMsg string
}

func (e *BadRequestError) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type NoSuchIDError struct {
	ErrorMsg string
}

func (e *NoSuchIDError) Error() string {
	return e.ErrorMsg
}
