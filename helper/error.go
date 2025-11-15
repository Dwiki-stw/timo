package helper

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) WriteError() {

}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

const (
	INTERNAL_ERROR string = "INTERNAL_SERVER_ERROR"
	INVALID_INPUT  string = "INVALID_INPUT"
	NOT_FOUND      string = "NOT_FOUND"
	LOGIN_ERROR    string = "LOGIN_ERROR"
	EMAIL_EXIST    string = "EMAIL_EXIST"
)
