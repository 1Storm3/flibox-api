package httperror

import (
	"errors"
	"fmt"
)

type Error struct {
	code    int
	message string
}

func New(code int, message string) error {
	return &Error{
		code:    code,
		message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Is(target error) bool {
	var t *Error
	ok := errors.As(target, &t)

	if !ok {
		return false
	}

	return e.code == t.code
}
