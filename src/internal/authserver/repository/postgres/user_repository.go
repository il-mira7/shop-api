package postgres

import (
	"context"
	"errors"
	"fmt"
	"shopAPI/internal/authserver/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	var uDAO userDAO
	uDAO.fromDomain(user)

	query := `
	INSERT INTO users (email, name, surname, phone_number, password_hash) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, created_at, updated_at;
	`
	err := r.db.QueryRow(ctx, query, uDAO.Email, uDAO.Name, uDAO.Surname, uDAO.PhoneNumber, uDAO.PasswordHash).Scan(&uDAO.ID, &uDAO.CreatedAt, &uDAO.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &domain.AlreadyExistsError{Entity: "user", Key: user.Email}
		}
		return fmt.Errorf("insert user: %w", err)
	}
	user.ID = uDAO.ID
	user.CreatedAt = uDAO.CreatedAt
	user.UpdatedAt = uDAO.UpdatedAt

	return nil
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `
	SELECT id, email, name, surname, phone_number, password_hash, created_at, updated_at
	FROM users
	WHERE LOWER(email) = LOWER($1)
	`

	rows, err := r.db.Query(ctx, query, login)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	uDAO, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userDAO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.NotFoundError{
				Entity: "user",
				ID:     login,
			}
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return uDAO.toDomain(), nil
}
func (r *userRepository) UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte) (bool, error) {
	query := `
	UPDATE users 
	SET password_hash = $2 
	WHERE id = $1;
	`
	result, err := r.db.Exec(ctx, query, id, passwordHash)
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return false, &domain.NotFoundError{
			Entity: "user",
			ID:     id.String(),
		}
	}
	return true, nil
}
