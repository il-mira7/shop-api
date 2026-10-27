package handler

import (
	"shopAPI/internal/shop/domain"

	"github.com/google/uuid"
)

type addressDTO struct {
	ID      uuid.UUID `json:"id,omitempty" format:"uuid" example:"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"`
	Country string    `json:"country" binding:"required" example:"Россия"`
	City    string    `json:"city" binding:"required" example:"Москва"`
	Street  string    `json:"street" binding:"required" example:"ул. Ленина, д. 10"`
}

func (a *addressDTO) toDomain() *domain.Address {
	return &domain.Address{
		ID:      a.ID,
		Country: a.Country,
		City:    a.City,
		Street:  a.Street,
	}
}

func (a *addressDTO) fromDomain(d *domain.Address) {
	a.ID = d.ID
	a.Country = d.Country
	a.City = d.City
	a.Street = d.Street
}
