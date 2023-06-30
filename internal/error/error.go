package error

import "fmt"

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

var (
	ErrEmailInUse              = &CustomError{Code: 400, Message: "Email is already in use"}
	ErrUserNotFound            = &CustomError{Code: 404, Message: "User not found"}
	ErrInvalidCredentials      = &CustomError{Code: 401, Message: "Invalid credentials"}
	ErrInternal                = &CustomError{Code: 500, Message: "Internal server error"}
	ErrSendingEmail            = &CustomError{Code: 500, Message: "Error in sending email"}
	ErrInvalidVerificationCode = &CustomError{Code: 400, Message: "Verification code is invalid"}
)
