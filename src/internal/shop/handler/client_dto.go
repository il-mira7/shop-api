package handler

import (
	"shopAPI/internal/shop/domain"
	"time"

	"github.com/google/uuid"
)

type clientDTO struct {
	ID               uuid.UUID  `json:"id,omitempty" format:"uuid" example:"d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44"`
	ClientName       string     `json:"name" binding:"required,min=2,max=100" example:"Иван"`
	ClientSurname    string     `json:"surname" binding:"required,min=2,max=100" example:"Иванов"`
	Birthday         string     `json:"birthday" binding:"required,datetime=2006-01-02" example:"1995-05-20"`
	Gender           string     `json:"gender" binding:"required,oneof=male female" example:"male"`
	RegistrationDate time.Time  `json:"-"`
	AddressID        *uuid.UUID `json:"-"`
	Address          addressDTO `json:"address" binding:"required"`
}

func (c *clientDTO) toDomain() *domain.Client {
	birthday, _ := time.Parse("2006-01-02", c.Birthday)
	address := c.Address.toDomain()
	return &domain.Client{
		ID:            c.ID,
		ClientName:    c.ClientName,
		ClientSurname: c.ClientSurname,
		Birthday:      birthday,
		Gender:        c.Gender,
		AddressID:     c.AddressID,
		Address:       address,
	}
}
func (c *clientDTO) fromDomain(d *domain.Client) {
	c.ID = d.ID
	c.ClientName = d.ClientName
	c.ClientSurname = d.ClientSurname
	c.Birthday = d.Birthday.Format("2006-01-02")
	c.Gender = d.Gender
	c.RegistrationDate = d.RegistrationDate
	c.AddressID = d.AddressID

	var aDTO addressDTO
	aDTO.fromDomain(d.Address)
	c.Address = aDTO

}

type clientsListResponse struct {
	Clients []clientDTO `json:"clients"`
	Total   int         `json:"total"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
}
