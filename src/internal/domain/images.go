package domain

import (
	"github.com/google/uuid"
)

type Images struct {
	ID    uuid.UUID
	Image []byte
}
