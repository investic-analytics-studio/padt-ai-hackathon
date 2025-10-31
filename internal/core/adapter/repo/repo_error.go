package repo

import (
	"fmt"
)

// ErrRecordNotFound is returned when a record cannot be found in the database
type ErrRecordNotFound struct {
	Table string
	ID    string
}

func (e *ErrRecordNotFound) Error() string {
	return fmt.Sprintf("record not found in %s table with id: %s", e.Table, e.ID)
}

// ErrDuplicateRecord is returned when trying to create a record that already exists
type ErrDuplicateRecord struct {
	Table string
	Field string
	Value string
}

func (e *ErrDuplicateRecord) Error() string {
	return fmt.Sprintf("duplicate record in %s table with %s: %s", e.Table, e.Field, e.Value)
}

// ErrInvalidOperation is returned when an invalid operation is attempted
type ErrInvalidOperation struct {
	Operation string
	Reason    string
}

func (e *ErrInvalidOperation) Error() string {
	return fmt.Sprintf("invalid operation %s: %s", e.Operation, e.Reason)
}

// ErrDatabaseOperation is returned when a database operation fails
type ErrDatabaseOperation struct {
	Operation string
	Err       error
}

func (e *ErrDatabaseOperation) Error() string {
	return fmt.Sprintf("database operation failed: %s - %v", e.Operation, e.Err)
}

func (e *ErrDatabaseOperation) Unwrap() error {
	return e.Err
}

// ErrRefcodeGeneration is returned when a refcode generation fails
type ErrRefcodeGeneration struct {
	Attempts int
	Err      error
}

func (e *ErrRefcodeGeneration) Error() string {
	return fmt.Sprintf("failed to generate refcode after %d attempts: %v", e.Attempts, e.Err)
}

func (e *ErrRefcodeGeneration) Unwrap() error {
	return e.Err
}
