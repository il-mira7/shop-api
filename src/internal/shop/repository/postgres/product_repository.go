package postgres

import (
	"context"
	"errors"
	"fmt"
	"shopAPI/internal/shop/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// foreignKeyViolationCode — код ошибки Postgres "foreign_key_violation"
// (https://www.postgresql.org/docs/current/errcodes-appendix.html).
const foreignKeyViolationCode = "23503"

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	product.LastUpdateDate = time.Now()

	var pDAO productDAO
	pDAO.fromDomain(product)

	query := `
	INSERT INTO product(name, category, price, available_stock, last_update_date, supplier_id, image_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id;
	`
	err := r.db.QueryRow(ctx, query, pDAO.Name, pDAO.Category, pDAO.Price, pDAO.AvailableStock, pDAO.LastUpdateDate, pDAO.SupplierID, pDAO.ImageID).Scan(&pDAO.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationCode {
			switch pgErr.ConstraintName {
			case "product_supplier_id_fkey":
				return &domain.ValidationError{Message: fmt.Sprintf("supplier with id %s not found", pDAO.SupplierID)}
			case "product_image_id_fkey":
				return &domain.ValidationError{Message: fmt.Sprintf("image with id %s not found", pDAO.ImageID)}
			}
		}
		return fmt.Errorf("insert product: %w", err)
	}
	product.ID = pDAO.ID

	return nil
}

func (r *productRepository) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	query := `
	UPDATE product
	SET available_stock = available_stock - $2, last_update_date = NOW()
	WHERE id = $1 AND available_stock >= $2;
	`
	result, err := r.db.Exec(ctx, query, id, count)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		var available int
		err := r.db.QueryRow(ctx, `SELECT available_stock FROM product WHERE id = $1;`, id).Scan(&available)
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.NotFoundError{
				Entity: "product",
				ID:     id.String(),
			}
		}
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}
		return &domain.ValidationError{
			Message: fmt.Sprintf("insufficient stock: available %d, requested %d", available, count),
		}
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
	SELECT id, name, category, price, available_stock, last_update_date, supplier_id, image_id
	FROM product
	WHERE id = $1;
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	pDAO, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[productDAO])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.NotFoundError{
				Entity: "product",
				ID:     id.String(),
			}
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return pDAO.toDomain(), nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	var total int
	queryCount := `
	SELECT count(*)
	FROM product
	WHERE available_stock > 0;
	`
	err := r.db.QueryRow(ctx, queryCount).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	query := `
	SELECT id, name, category, price, available_stock, last_update_date, supplier_id, image_id
	FROM product
    WHERE available_stock > 0
    ORDER BY last_update_date DESC
    LIMIT (CASE WHEN $1 > 0 THEN $1 ELSE NULL END)
    OFFSET (CASE WHEN $2 > 0 THEN $2 ELSE 0 END);
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}

	pDAO, err := pgx.CollectRows(rows, pgx.RowToStructByName[productDAO])
	if err != nil {
		return nil, 0, fmt.Errorf("scan failed: %w", err)
	}

	product := make([]*domain.Product, 0, len(pDAO))
	for _, v := range pDAO {
		product = append(product, v.toDomain())
	}

	return product, total, nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM product
	WHERE id = $1;
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return &domain.NotFoundError{
			Entity: "product",
			ID:     id.String(),
		}
	}

	return nil
}
