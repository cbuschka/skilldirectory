package errors

// Implements error interface
type MissingIDError struct {
	ErrorMsg string
}

func (e *MissingIDError) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type NoSuchIDError struct {
	ErrorMsg string
}

func (e *NoSuchIDError) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type InvalidSkillTypeError struct {
	ErrorMsg string
}

func (e *InvalidSkillTypeError) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type MarshalingError struct {
	ErrorMsg string
}

func (e *MarshalingError) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type SavingError struct {
	ErrorMsg string
}

func (e *SavingError) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type IncompletePOSTBodyError struct {
	ErrorMsg string
}

func (e *IncompletePOSTBodyError) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type InvalidPOSTBodyError struct {
	ErrorMsg string
}

func (e *InvalidPOSTBodyError) Error() string {
	return e.ErrorMsg
}
