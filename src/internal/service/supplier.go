package service

import (
	"context"
	"shopAPI/internal/domain"
	"strings"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *domain.Supplier) error
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Supplier, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) error
}

type supplierService struct {
	repo SupplierRepository
}

func NewSupplierService(repo SupplierRepository) *supplierService {
	return &supplierService{repo: repo}
}

func (s *supplierService) Create(ctx context.Context, supplier *domain.Supplier) error {
	if supplier == nil {
		return &domain.EmptyPayloadError{Entity: "supplier"}
	}

	supplier.PhoneNumber = cleanPhoneNumber(supplier.PhoneNumber)
	if !isValidPhoneNumber(supplier.PhoneNumber) {
		return &domain.ValidationError{Message: "invalid phone format"}
	}

	var ok bool
	supplier.Name, ok = checkName(supplier.Name)
	if !ok {
		return &domain.ValidationError{Message: "name is too short or invalid"}
	}

	if err := checkAddress(supplier.Address); err != nil {
		return err
	}

	return s.repo.Create(ctx, supplier)
}

func (s *supplierService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Supplier, int, error) {
	if err := validatePagination(limit, offset); err != nil {
		return nil, 0, err
	}
	return s.repo.GetAll(ctx, limit, offset)
}
func (s *supplierService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *supplierService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *supplierService) Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) (*domain.Supplier, error) {
	if err := checkAddress(newAddress); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, id, newAddress); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func cleanPhoneNumber(phone string) string {
	phone = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '+' {
			return r
		}
		return -1
	}, phone)

	if strings.HasPrefix(phone, "8") {
		phone = "+7" + strings.TrimPrefix(phone, "8")
	}
	return phone
}

func isValidPhoneNumber(phone string) bool {
	ph, err := phonenumbers.Parse(phone, "RU")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(ph)
}

func checkName(name string) (string, bool) {
	words := strings.Fields(name)

	name = strings.Join(words, " ")
	if len(name) < 2 {
		return name, false
	}
	return name, true
}
