package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID
	Name           string
	Category       string
	Price          float64
	AvailableStock uint
	LastUpdateDate time.Time
	SupplierID     *uuid.UUID
	ImageID        *uuid.UUID
}
