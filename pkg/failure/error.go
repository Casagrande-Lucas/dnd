package failure

import "errors"

var (
	ErrorBadRequest             = errors.New("bad request")
	ErrorUnauthorized           = errors.New("unauthorized")
	ErrorForbidden              = errors.New("forbidden")
	ErrorNotFound               = errors.New("not found")
	ErrorMethodNotAllowed       = errors.New("method not allowed")
	ErrorNotAcceptable          = errors.New("not acceptable")
	ErrorInternalServer         = errors.New("internal server error")
	ErrorDeadlineExceeded       = errors.New("deadline exceeded")
	ErrorEmailAlreadyRegistered = errors.New("email already registered")
	ErrorMigrate                = errors.New("migrate filed")
)

type Error struct {
	appErr error
	svcErr error
}

// NewError creates a new custom Error instance
func NewError(app error, svc error) *Error {
	return &Error{appErr: app, svcErr: svc}
}

// Error returns the combined error message
func (e Error) Error() string {
	joinedErr := errors.Join(e.svcErr, e.appErr)
	if joinedErr == nil {
		return ""
	}
	return joinedErr.Error()
}

// AppErr retrieves the application error
func (e Error) AppErr() error {
	return e.appErr
}

// SvcErr retrieves the service error
func (e Error) SvcErr() error {
	return e.svcErr
}

// Unwrap allows errors.Is and errors.As to traverse wrapped errors
func (e Error) Unwrap() []error {
	return []error{e.svcErr, e.appErr}
}
