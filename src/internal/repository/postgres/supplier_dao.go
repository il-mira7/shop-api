package postgres

import (
	"shopAPI/internal/domain"

	"github.com/google/uuid"
)

type supplierDAO struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	AddressID   *uuid.UUID `db:"address_id"`
	PhoneNumber string     `db:"phone_number"`
}

func (s *supplierDAO) toDomain(address *domain.Address) *domain.Supplier {
	return &domain.Supplier{
		ID:          s.ID,
		Name:        s.Name,
		AddressID:   s.AddressID,
		Address:     address,
		PhoneNumber: s.PhoneNumber,
	}
}

func (s *supplierDAO) fromDomain(d *domain.Supplier) {
	s.ID = d.ID
	s.Name = d.Name
	s.AddressID = d.AddressID
	s.PhoneNumber = d.PhoneNumber
}
