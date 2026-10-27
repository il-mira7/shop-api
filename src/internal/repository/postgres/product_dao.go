package postgres

import (
	"shopAPI/internal/domain"
	"time"

	"github.com/google/uuid"
)

type productDAO struct {
	ID             uuid.UUID  `db:"id"`
	Name           string     `db:"name"`
	Category       string     `db:"category"`
	Price          float64    `db:"price"`
	AvailableStock uint       `db:"available_stock"`
	LastUpdateDate time.Time  `db:"last_update_date"`
	SupplierID     *uuid.UUID `db:"supplier_id"`
	ImageID        *uuid.UUID `db:"image_id"`
}

func (p *productDAO) toDomain() *domain.Product {
	return &domain.Product{
		ID:             p.ID,
		Name:           p.Name,
		Category:       p.Category,
		Price:          p.Price,
		AvailableStock: p.AvailableStock,
		LastUpdateDate: p.LastUpdateDate,
		SupplierID:     p.SupplierID,
		ImageID:        p.ImageID,
	}
}

func (p *productDAO) fromDomain(d *domain.Product) {
	p.ID = d.ID
	p.Name = d.Name
	p.Category = d.Category
	p.Price = d.Price
	p.AvailableStock = d.AvailableStock
	p.LastUpdateDate = d.LastUpdateDate
	p.SupplierID = d.SupplierID
	p.ImageID = d.ImageID
}
