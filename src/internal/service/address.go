package service

import (
	"shopAPI/internal/domain"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func checkAddress(address *domain.Address) error {
	if address == nil {
		return &domain.EmptyPayloadError{Entity: "address"}
	}
	caser := cases.Title(language.Russian)
	address.City = caser.String(strings.ToLower(strings.TrimSpace(address.City)))
	address.Country = caser.String(strings.ToLower(strings.TrimSpace(address.Country)))
	address.Street = strings.TrimSpace(address.Street)

	if address.City == "" || address.Country == "" || address.Street == "" {
		return &domain.ValidationError{Message: "all address fields (city, country, street) must be filled"}
	}
	return nil
}
