package sdk

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInternalServer = errors.New("internal error")
	ErrForbidden      = errors.New("forbidden")
	ErrUnauthorized   = errors.New("unauthorized")
)

type TooManyRequestsError struct {
	RetryAfter time.Duration
}

func NewTooManyRequestsError(retryAfter time.Duration) TooManyRequestsError {
	return TooManyRequestsError{
		RetryAfter: retryAfter,
	}
}

func (e TooManyRequestsError) Error() string {
	return fmt.Sprintf("too many requests [retry after: %.2f]", e.RetryAfter.Seconds())
}

type ValidationError struct {
	msg    string
	errors map[string]interface{}
}

func NewValidationError(msg string, errs map[string]interface{}) ValidationError {
	return ValidationError{
		msg:    msg,
		errors: errs,
	}
}

func (e ValidationError) Error() string {
	return e.msg
}
