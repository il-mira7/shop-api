package postgres

import (
	"context"
	"errors"
	"fmt"
	"shopAPI/internal/shop/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type clientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *clientRepository {
	return &clientRepository{db: db}
}

type clientResultDAO struct {
	clientDAO
	addressDAO
}

func (r *clientRepository) Create(ctx context.Context, client *domain.Client) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	addressID, err := getOrCreateAddress(ctx, tx, client.Address)
	if err != nil {
		return fmt.Errorf("insert address: %w", err)
	}

	var cDAO clientDAO
	client.RegistrationDate = time.Now()
	cDAO.fromDomain(client)
	query := `
	INSERT INTO client (client_name, client_surname, birthday, gender, registration_date, address_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`
	err = tx.QueryRow(ctx, query, cDAO.ClientName, cDAO.ClientSurname, cDAO.Birthday, cDAO.Gender, cDAO.RegistrationDate, addressID).Scan(&cDAO.ID)
	if err != nil {
		return fmt.Errorf("insert client: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	client.ID = cDAO.ID
	client.AddressID = &addressID
	if client.Address != nil {
		client.Address.ID = addressID
	}

	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE 
	FROM client
	WHERE id = $1;
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return &domain.NotFoundError{
			Entity: "client",
			ID:     id.String(),
		}
	}

	return nil
}

func (r *clientRepository) GetByNames(ctx context.Context, firstName string, lastName string) ([]*domain.Client, error) {
	query :=
		`
		SELECT c.id, c.client_name, c.client_surname, c.birthday, c.gender, c.registration_date, c.address_id, a.country, a.city, a.street
		FROM client c
		JOIN address a ON c.address_id = a.id
		WHERE c.client_name = $1 AND c.client_surname = $2
		ORDER BY c.registration_date;
		`
	rows, err := r.db.Query(ctx, query, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	cDAO, err := pgx.CollectRows(rows, pgx.RowToStructByName[clientResultDAO])
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	clients := make([]*domain.Client, len(cDAO))
	for i, v := range cDAO {
		address := v.addressDAO.toDomain()
		if v.AddressID != nil {
			address.ID = *v.AddressID
		}
		clients[i] = v.clientDAO.toDomain(address)
	}
	return clients, nil
}

func (r *clientRepository) GetAll(ctx context.Context, limit int, offset int) ([]*domain.Client, int, error) {
	queryCount :=
		`
		SELECT count(*)
		FROM client;
		`
	var total int
	err := r.db.QueryRow(ctx, queryCount).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	query :=
		`
		SELECT c.id, c.client_name, c.client_surname, c.birthday, c.gender, c.registration_date, c.address_id, a.country, a.city, a.street
		FROM client c
		JOIN address a ON c.address_id = a.id
		ORDER BY c.id
		LIMIT (CASE WHEN $1 > 0 THEN $1 ELSE NULL END)
		OFFSET (CASE WHEN $2 > 0 THEN $2 ELSE 0 END);
		`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	cDAO, err := pgx.CollectRows(rows, pgx.RowToStructByName[clientResultDAO])
	if err != nil {
		return nil, 0, fmt.Errorf("scan failed: %w", err)
	}
	clients := make([]*domain.Client, len(cDAO))
	for i, v := range cDAO {
		address := v.addressDAO.toDomain()
		if v.AddressID != nil {
			address.ID = *v.AddressID
		}
		clients[i] = v.clientDAO.toDomain(address)
	}
	return clients, total, nil
}

func (r *clientRepository) Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	newAddressID, err := getOrCreateAddress(ctx, tx, newAddress)
	if err != nil {
		return fmt.Errorf("get or create address: %w", err)
	}

	query :=
		`
		UPDATE client
		SET address_id = $2
		WHERE id = $1;
		`
	result, err := tx.Exec(ctx, query, id, newAddressID)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return &domain.NotFoundError{
			Entity: "client",
			ID:     id.String(),
		}
	}

	return tx.Commit(ctx)
}

func (r *clientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	query :=
		`
		SELECT c.id, c.client_name, c.client_surname, c.birthday, c.gender, c.registration_date, c.address_id, a.country, a.city, a.street
		FROM client c
		JOIN address a ON c.address_id = a.id
		WHERE c.id = $1;
		`

	row, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	result, err := pgx.CollectOneRow(row, pgx.RowToStructByName[clientResultDAO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.NotFoundError{
				Entity: "client",
				ID:     id.String(),
			}
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	address := result.addressDAO.toDomain()
	if result.AddressID != nil {
		address.ID = *result.AddressID
	}
	client := result.clientDAO.toDomain(address)

	return client, nil
}
