package customerror

import "fmt"

type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrEmailInUse              = &CustomError{Code: "EMAIL_IN_USE", Message: "Email is already in use"}
	ErrUserNotFound            = &CustomError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrInvalidCredentials      = &CustomError{Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
	ErrInternal                = &CustomError{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error"}
	ErrSendingEmail            = &CustomError{Code: "ERROR_SENDING_EMAIL", Message: "Error in sending email"}
	ErrInvalidVerificationCode = &CustomError{Code: "INVALID_VERIFICATION_CODE", Message: "Verification code is invalid"}
	ErrBadRequest              = &CustomError{Code: "BAD_REQUEST", Message: "Request is invalid"}
	ErrUnauthorized            = &CustomError{Code: "UNAUTHORIZED", Message: "Token is unauthorized"}
)
