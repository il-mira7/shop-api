package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID               uuid.UUID
	ClientName       string
	ClientSurname    string
	Birthday         time.Time
	Gender           string
	RegistrationDate time.Time
	AddressID        *uuid.UUID
	Address          *Address
}
