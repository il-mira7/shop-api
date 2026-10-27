package postgres

import (
	"shopAPI/internal/domain"
	"time"

	"github.com/google/uuid"
)

type clientDAO struct {
	ID               uuid.UUID  `db:"id"`
	ClientName       string     `db:"client_name"`
	ClientSurname    string     `db:"client_surname"`
	Birthday         time.Time  `db:"birthday"`
	Gender           string     `db:"gender"`
	RegistrationDate time.Time  `db:"registration_date"`
	AddressID        *uuid.UUID `db:"address_id"`
}

func (c *clientDAO) toDomain(address *domain.Address) *domain.Client {
	return &domain.Client{
		ID:               c.ID,
		ClientName:       c.ClientName,
		ClientSurname:    c.ClientSurname,
		Birthday:         c.Birthday,
		Gender:           c.Gender,
		RegistrationDate: c.RegistrationDate,
		AddressID:        c.AddressID,
		Address:          address,
	}
}

func (c *clientDAO) fromDomain(d *domain.Client) {
	c.ID = d.ID
	c.ClientName = d.ClientName
	c.ClientSurname = d.ClientSurname
	c.Birthday = d.Birthday
	c.Gender = d.Gender
	c.RegistrationDate = d.RegistrationDate
	c.AddressID = d.AddressID
}
