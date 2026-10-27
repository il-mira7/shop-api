package service

import (
	"context"
	"shopAPI/internal/domain"

	"github.com/google/uuid"
)

const maxFileSize = 5 * 1024 * 1024

type ImageRepository interface {
	Create(ctx context.Context, image []byte, productID uuid.UUID) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, image []byte) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]byte, error)
	GetByImageID(ctx context.Context, id uuid.UUID) ([]byte, error)
}

type imageService struct {
	repo ImageRepository
}

func NewImageService(repo ImageRepository) *imageService {
	return &imageService{repo: repo}
}

func (s *imageService) Create(ctx context.Context, image []byte, productID uuid.UUID) (uuid.UUID, error) {
	if productID == uuid.Nil {
		return uuid.Nil, &domain.ValidationError{Message: "product id is required"}
	}

	if len(image) == 0 {
		return uuid.Nil, &domain.ValidationError{Message: "image data is empty"}
	}

	if len(image) > maxFileSize {
		return uuid.Nil, &domain.ValidationError{Message: "file is too large"}
	}
	return s.repo.Create(ctx, image, productID)
}
func (s *imageService) Update(ctx context.Context, id uuid.UUID, image []byte) error {
	if id == uuid.Nil {
		return &domain.ValidationError{Message: "id is required"}
	}
	if len(image) == 0 {
		return &domain.ValidationError{Message: "image data is empty"}
	}

	if len(image) > maxFileSize {
		return &domain.ValidationError{Message: "file is too large"}
	}
	return s.repo.Update(ctx, id, image)
}
func (s *imageService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return &domain.ValidationError{Message: "id is required"}
	}
	return s.repo.Delete(ctx, id)
}
func (s *imageService) GetByProductID(ctx context.Context, productID uuid.UUID) ([]byte, error) {
	if productID == uuid.Nil {
		return nil, &domain.ValidationError{Message: "product id is required"}
	}
	return s.repo.GetByProductID(ctx, productID)
}
func (s *imageService) GetByImageID(ctx context.Context, id uuid.UUID) ([]byte, error) {
	if id == uuid.Nil {
		return nil, &domain.ValidationError{Message: "id is required"}
	}
	return s.repo.GetByImageID(ctx, id)
}
