package exceptions

import "fmt"

type AppError struct {
	StatusCode int
	Message    string
	Cause      error
}

func NewAppError(statusCode int, message string, cause error) *AppError {
	return &AppError{StatusCode: statusCode, Message: message, Cause: cause}
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprintf("application error with status %d", e.StatusCode)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}
