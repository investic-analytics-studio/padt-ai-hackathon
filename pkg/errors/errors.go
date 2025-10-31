package errors

import "fmt"

type AppError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error_code"`
}

func (e AppError) Error() string { return fmt.Sprintf("[%v]: %v", e.ErrorCode, e.Message) }
