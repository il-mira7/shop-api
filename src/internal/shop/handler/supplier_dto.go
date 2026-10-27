package handler

import (
	"shopAPI/internal/shop/domain"

	"github.com/google/uuid"
)

type supplierDTO struct {
	ID          uuid.UUID  `json:"id,omitempty" format:"uuid" example:"b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22"`
	Name        string     `json:"name" binding:"required" example:"ООО ТехноПарк"`
	AddressID   *uuid.UUID `json:"-"`
	Address     addressDTO `json:"address" binding:"required"`
	PhoneNumber string     `json:"phone_number" binding:"required" example:"+79991234567"`
}

func (s *supplierDTO) toDomain() *domain.Supplier {
	address := s.Address.toDomain()
	return &domain.Supplier{
		ID:          s.ID,
		Name:        s.Name,
		AddressID:   s.AddressID,
		Address:     address,
		PhoneNumber: s.PhoneNumber,
	}
}

func (s *supplierDTO) fromDomain(d *domain.Supplier) {
	s.ID = d.ID
	s.Name = d.Name
	s.AddressID = d.AddressID

	var aDTO addressDTO
	aDTO.fromDomain(d.Address)
	s.Address = aDTO

	s.PhoneNumber = d.PhoneNumber
}

type suppliersListResponse struct {
	Suppliers []supplierDTO `json:"suppliers"`
	Total     int           `json:"total"`
	Limit     int           `json:"limit"`
	Offset    int           `json:"offset"`
}
