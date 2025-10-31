package errors

import "net/http"

func New(status int, code, message string) AppError {
	return AppError{
		ErrorCode:  code,
		StatusCode: status,
		Message:    message,
	}
}

func NewBadRequest(code, message string) AppError {
	return New(http.StatusBadRequest, code, message)
}

func NewNotFound(code, message string) AppError {
	return New(http.StatusNotFound, code, message)
}

func NewUnprocessable(code, message string) AppError {
	return New(http.StatusUnprocessableEntity, code, message)
}

func NewUnauthorized(code, message string) AppError {
	return New(http.StatusUnauthorized, code, message)
}

func NewForbidden(code, message string) AppError {
	return New(http.StatusForbidden, code, message)
}

func NewConflict(code, message string) AppError {
	return New(http.StatusConflict, code, message)
}

func NewInternalServer(code, message string) AppError {
	return New(http.StatusInternalServerError, code, message)
}
