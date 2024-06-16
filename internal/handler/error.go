package handler

type HandlerError struct {
	code int
	msg  string
}

func (e HandlerError) Error() string {
	return e.msg
}

func NewHandlerError(code int, msg string) *HandlerError {
	return &HandlerError{
		code: code,
		msg:  msg,
	}
}

var (
	ErrNotFound   = NewHandlerError(404, "not found")
	ErrConflict   = NewHandlerError(409, "conflict")
	ErrBadRequest = NewHandlerError(400, "bad request")
)
