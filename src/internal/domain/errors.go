package domain

import "fmt"

type NotFoundError struct {
	Entity string
	ID     string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.Entity, e.ID)
}

type ValidationError struct {
	Message string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", ve.Message)
}

type EmptyPayloadError struct {
	Entity string
}

func (e *EmptyPayloadError) Error() string {
	return fmt.Sprintf("the object %s is null", e.Entity)
}
