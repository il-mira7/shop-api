package service

import "shopAPI/internal/domain"

const maxLimit = 100

func validatePagination(limit, offset int) error {
	if limit < 0 || offset < 0 {
		return &domain.ValidationError{Message: "limit or offset cannot be negative"}
	}
	if limit > maxLimit {
		return &domain.ValidationError{Message: "limit is too big"}
	}
	return nil
}
