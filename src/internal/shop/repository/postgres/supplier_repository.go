package postgres

import (
	"context"
	"errors"
	"fmt"
	"shopAPI/internal/shop/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type supplierRepository struct {
	db *pgxpool.Pool
}

func NewSupplierRepository(db *pgxpool.Pool) *supplierRepository {
	return &supplierRepository{db: db}
}

type resultDAO struct {
	supplierDAO
	addressDAO
}

func (r *supplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	addressID, err := getOrCreateAddress(ctx, tx, supplier.Address)
	if err != nil {
		return fmt.Errorf("insert address: %w", err)
	}

	var sDAO supplierDAO
	sDAO.fromDomain(supplier)

	query := `
	INSERT INTO supplier(name, address_id, phone_number)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	err = tx.QueryRow(ctx, query, sDAO.Name, addressID, sDAO.PhoneNumber).Scan(&sDAO.ID)
	if err != nil {
		return fmt.Errorf("insert supplier: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	supplier.ID = sDAO.ID
	supplier.AddressID = &addressID
	if supplier.Address != nil {
		supplier.Address.ID = addressID
	}

	return nil
}

func (r *supplierRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Supplier, int, error) {
	var total int
	queryCount := `
	SELECT COUNT(*)
	FROM supplier;
	`
	err := r.db.QueryRow(ctx, queryCount).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	query := `
		SELECT s.id, s.name, s.address_id, s.phone_number, a.country, a.city, a.street
		FROM supplier s
		JOIN address a ON s.address_id = a.id
		ORDER BY s.id
    LIMIT (CASE WHEN $1 > 0 THEN $1 ELSE NULL END)
    OFFSET (CASE WHEN $2 > 0 THEN $2 ELSE 0 END);
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[resultDAO])
	if err != nil {
		return nil, 0, fmt.Errorf("scan failed: %w", err)
	}

	suppliers := make([]*domain.Supplier, 0, len(result))
	for _, v := range result {
		address := v.addressDAO.toDomain()
		if v.AddressID != nil {
			address.ID = *v.AddressID
		}
		suppliers = append(suppliers, v.supplierDAO.toDomain(address))
	}
	return suppliers, total, nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	query := `
		SELECT s.id, s.name, s.address_id, s.phone_number, a.country, a.city, a.street
		FROM supplier s
		JOIN address a ON s.address_id = a.id
		WHERE s.id = $1;
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[resultDAO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.NotFoundError{
				Entity: "supplier",
				ID:     id.String(),
			}
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	address := result.addressDAO.toDomain()
	if result.AddressID != nil {
		address.ID = *result.AddressID
	}
	supplier := result.supplierDAO.toDomain(address)
	return supplier, nil
}

func (r *supplierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM supplier
	WHERE id = $1; 
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return &domain.NotFoundError{
			Entity: "supplier",
			ID:     id.String(),
		}
	}

	return nil
}

func (r *supplierRepository) Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	newAddressID, err := getOrCreateAddress(ctx, tx, newAddress)
	if err != nil {
		return fmt.Errorf("get or create address: %w", err)
	}

	query := `
	UPDATE supplier
	SET address_id = $2
	WHERE id = $1;
	`

	result, err := tx.Exec(ctx, query, id, newAddressID)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return &domain.NotFoundError{
			Entity: "supplier",
			ID:     id.String(),
		}
	}

	return tx.Commit(ctx)
}
