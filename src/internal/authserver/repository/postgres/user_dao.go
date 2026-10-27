package postgres

import (
	"shopAPI/internal/authserver/domain"
	"time"

	"github.com/google/uuid"
)

type userDAO struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	Surname      string    `db:"surname"`
	PhoneNumber  string    `db:"phone_number"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (u *userDAO) toDomain() *domain.User {
	return &domain.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Surname:      u.Surname,
		PhoneNumber:  u.PhoneNumber,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (u *userDAO) fromDomain(d *domain.User) {
	u.Email = d.Email
	u.Name = d.Name
	u.Surname = d.Surname
	u.PhoneNumber = d.PhoneNumber
	u.PasswordHash = d.PasswordHash
}
