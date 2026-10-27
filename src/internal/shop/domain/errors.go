package domain

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	Entity string
	ID     string
	Msg    string
}

func (e *NotFoundError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
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

type AlreadyExistsError struct {
	Entity string
	Key    string
	Msg    string
}

func (e *AlreadyExistsError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return fmt.Sprintf("%s %s already exists", e.Entity, e.Key)
}

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrAuthServiceUnavailable = errors.New("authorization service is unavailable")
