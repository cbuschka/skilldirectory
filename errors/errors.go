package errors

// Implements error interface
type BadRequest400Error struct {
	ErrorMsg string
}

func (e *BadRequest400Error) Error() string {
	return e.ErrorMsg
}

// Implements error interface
type NoSuchID404Error struct {
	ErrorMsg string
}

func (e *NoSuchID404Error) Error() string {
	return e.ErrorMsg
}
