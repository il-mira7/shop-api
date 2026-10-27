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

type imageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) *imageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) Create(ctx context.Context, image []byte, productID uuid.UUID) (uuid.UUID, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
	INSERT INTO images(image)
	VALUES ($1)
	RETURNING id;
	`
	var id uuid.UUID
	err = tx.QueryRow(ctx, query, image).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert image: %w", err)
	}

	queryUpdate := `
	UPDATE product
	SET image_id = $2
	WHERE id = $1;
	`
	result, err := tx.Exec(ctx, queryUpdate, productID, id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("update query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return uuid.Nil, &domain.NotFoundError{
			Entity: "product",
			ID:     productID.String(),
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit transaction: %w", err)
	}
	return id, nil
}

func (r *imageRepository) Update(ctx context.Context, id uuid.UUID, image []byte) error {
	query := `
	UPDATE images
	SET image = $2
	WHERE id = $1;
	`
	result, err := r.db.Exec(ctx, query, id, image)
	if err != nil {
		return fmt.Errorf("update query failed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return &domain.NotFoundError{
			Entity: "image",
			ID:     id.String(),
		}
	}
	return nil
}

func (r *imageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM images
	WHERE id = $1;
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return &domain.NotFoundError{
			Entity: "image",
			ID:     id.String(),
		}
	}
	return nil
}

func (r *imageRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]byte, error) {
	query := `
	SELECT image 
	FROM images i
	JOIN product p ON p.image_id = i.id
	WHERE p.id = $1;
	`
	var image []byte
	err := r.db.QueryRow(ctx, query, productID).Scan(&image)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.NotFoundError{
				Entity: "image on product ID",
				ID:     productID.String(),
			}
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return image, nil
}

func (r *imageRepository) GetByImageID(ctx context.Context, id uuid.UUID) ([]byte, error) {
	query := `
	SELECT image 
	FROM images
	WHERE id = $1;
	`
	var image []byte
	err := r.db.QueryRow(ctx, query, id).Scan(&image)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.NotFoundError{
				Entity: "image",
				ID:     id.String(),
			}
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return image, nil
}
