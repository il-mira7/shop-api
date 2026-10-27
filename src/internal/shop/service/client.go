package service

import (
	"context"
	"shopAPI/internal/shop/domain"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Create(ctx context.Context, client *domain.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByNames(ctx context.Context, firstName string, lastName string) ([]*domain.Client, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*domain.Client, int, error)
	Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error)
}

type clientService struct {
	repo ClientRepository
}

func NewClientService(repo ClientRepository) *clientService {
	return &clientService{repo: repo}
}

func (s *clientService) Create(ctx context.Context, client *domain.Client) error {
	if client == nil {
		return &domain.EmptyPayloadError{Entity: "client"}
	}
	if err := checkAddress(client.Address); err != nil {
		return err
	}
	return s.repo.Create(ctx, client)
}

func (s *clientService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *clientService) GetByNames(ctx context.Context, firstName string, lastName string) ([]*domain.Client, error) {
	return s.repo.GetByNames(ctx, firstName, lastName)
}

func (s *clientService) GetAll(ctx context.Context, limit int, offset int) ([]*domain.Client, int, error) {
	if err := validatePagination(limit, offset); err != nil {
		return nil, 0, err
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *clientService) Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) (*domain.Client, error) {
	if err := checkAddress(newAddress); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, id, newAddress); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}
