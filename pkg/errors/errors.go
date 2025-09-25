package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application-specific error codes
type ErrorCode string

const (
	// Domain error codes
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"

	// User domain errors
	ErrCodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeEmailExists        ErrorCode = "EMAIL_EXISTS"
	ErrCodeUserNotVerified    ErrorCode = "USER_NOT_VERIFIED"

	// Listing domain errors
	ErrCodeListingNotFound   ErrorCode = "LISTING_NOT_FOUND"
	ErrCodeListingInactive   ErrorCode = "LISTING_INACTIVE"
	ErrCodeInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"

	// Transaction domain errors
	ErrCodeTransactionNotFound ErrorCode = "TRANSACTION_NOT_FOUND"
	ErrCodePaymentFailed       ErrorCode = "PAYMENT_FAILED"
	ErrCodeEscrowError         ErrorCode = "ESCROW_ERROR"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *DomainError) WithDetails(key string, value interface{}) *DomainError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// HTTPStatusCode returns the appropriate HTTP status code for the error
func (e *DomainError) HTTPStatusCode() int {
	switch e.Code {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeNotFound, ErrCodeUserNotFound, ErrCodeListingNotFound, ErrCodeTransactionNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized, ErrCodeInvalidCredentials:
		return http.StatusUnauthorized
	case ErrCodeForbidden, ErrCodeUserNotVerified:
		return http.StatusForbidden
	case ErrCodeConflict, ErrCodeEmailExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Common error constructors
func ValidationError(message string) *DomainError {
	return NewDomainError(ErrCodeValidation, message)
}

func NotFoundError(message string) *DomainError {
	return NewDomainError(ErrCodeNotFound, message)
}

func UnauthorizedError(message string) *DomainError {
	return NewDomainError(ErrCodeUnauthorized, message)
}

func ConflictError(message string) *DomainError {
	return NewDomainError(ErrCodeConflict, message)
}

func InternalError(message string) *DomainError {
	return NewDomainError(ErrCodeInternalServer, message)
}
