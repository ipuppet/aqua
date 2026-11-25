package logic

import "net/http"

type AuthError struct {
	Message    string
	StatusCode int
}

func (e *AuthError) Error() string {
	return e.Message
}

func (e *AuthError) HTTPStatus() int {
	return e.StatusCode
}

func (e *AuthError) ClientMessage() string {
	return e.Message
}

func NewAuthError() *AuthError {
	return &AuthError{
		Message:    "account not found or password incorrect",
		StatusCode: http.StatusBadRequest,
	}
}

type AccountLockedError struct {
	Message    string
	StatusCode int
}

func (e *AccountLockedError) Error() string {
	return e.Message
}

func (e *AccountLockedError) HTTPStatus() int {
	return e.StatusCode
}

func (e *AccountLockedError) ClientMessage() string {
	return e.Message
}

func NewAccountLockedError() *AccountLockedError {
	return &AccountLockedError{
		Message:    "account has been locked",
		StatusCode: http.StatusForbidden,
	}
}

type TokenError struct {
	Message    string
	StatusCode int
}

func (e *TokenError) Error() string {
	return e.Message
}

func (e *TokenError) HTTPStatus() int {
	return e.StatusCode
}

func (e *TokenError) ClientMessage() string {
	return e.Message
}

func NewTokenError(message string) *TokenError {
	return &TokenError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

type ValidationError struct {
	Message    string
	StatusCode int
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) HTTPStatus() int {
	return e.StatusCode
}

func (e *ValidationError) ClientMessage() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
