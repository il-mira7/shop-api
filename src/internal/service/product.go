package service

import (
	"context"
	"shopAPI/internal/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) Create(ctx context.Context, product *domain.Product) error {
	if product == nil {
		return &domain.EmptyPayloadError{Entity: "product"}
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) DecreaseStock(ctx context.Context, id uuid.UUID, count int) error {
	if count <= 0 {
		return &domain.ValidationError{Message: "count is too small"}
	}
	return s.repo.DecreaseStock(ctx, id, count)
}

func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	if err := validatePagination(limit, offset); err != nil {
		return nil, 0, err
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
