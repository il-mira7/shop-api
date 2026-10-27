package postgres

import (
	"context"
	"errors"
	"shopAPI/internal/shop/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func createWithDB(ctx context.Context, db DBTX, address *domain.Address) (uuid.UUID, error) {
	var aDAO addressDAO
	aDAO.fromDomain(address)

	query := `
	INSERT INTO address(country, city, street)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	var newID uuid.UUID
	err := db.QueryRow(ctx, query, aDAO.Country, aDAO.City, aDAO.Street).Scan(&newID)
	return newID, err
}

func getOrCreateAddress(ctx context.Context, db DBTX, address *domain.Address) (uuid.UUID, error) {
	if address == nil {
		return uuid.Nil, &domain.EmptyPayloadError{Entity: "address"}
	}

	var aDAO addressDAO
	aDAO.fromDomain(address)

	query := `
	SELECT id
	FROM address
	WHERE country = $1 AND city = $2 AND street = $3;
	`
	var ID uuid.UUID
	err := db.QueryRow(ctx, query, aDAO.Country, aDAO.City, aDAO.Street).Scan(&ID)

	if errors.Is(err, pgx.ErrNoRows) {
		return createWithDB(ctx, db, address)
	}
	if err != nil {
		return uuid.Nil, err
	}

	return ID, nil

}
