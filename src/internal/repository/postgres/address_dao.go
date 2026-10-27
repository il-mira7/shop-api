package postgres

import (
	"context"
	"shopAPI/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type addressDAO struct {
	ID      uuid.UUID `db:"-"`
	Country string    `db:"country"`
	City    string    `db:"city"`
	Street  string    `db:"street"`
}

func (a *addressDAO) toDomain() *domain.Address {
	return &domain.Address{
		ID:      a.ID,
		Country: a.Country,
		City:    a.City,
		Street:  a.Street,
	}
}

func (a *addressDAO) fromDomain(d *domain.Address) {
	a.ID = d.ID
	a.Country = d.Country
	a.City = d.City
	a.Street = d.Street
}

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
