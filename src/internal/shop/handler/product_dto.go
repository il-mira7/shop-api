package handler

import (
	"shopAPI/internal/shop/domain"
	"time"

	"github.com/google/uuid"
)

type productDTO struct {
	ID             uuid.UUID  `json:"id,omitempty" format:"uuid" example:"c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33"`
	Name           string     `json:"name" binding:"required" example:"Холодильник"`
	Category       string     `json:"category" binding:"required" example:"Бытовая техника"`
	Price          float64    `json:"price" binding:"required" example:"49990.50"`
	AvailableStock uint       `json:"available_stock" binding:"required" example:"10"`
	LastUpdateDate time.Time  `json:"-"`
	SupplierID     *uuid.UUID `json:"supplier_id,omitempty" format:"uuid"`
	ImageID        *uuid.UUID `json:"image_id,omitempty" format:"uuid"`
}

type decreaseStockRequest struct {
	Count int `json:"count" binding:"required,gt=0" example:"2"`
}

func (p *productDTO) toDomain() *domain.Product {
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

func (p *productDTO) fromDomain(d *domain.Product) {
	p.ID = d.ID
	p.Name = d.Name
	p.Category = d.Category
	p.Price = d.Price
	p.AvailableStock = d.AvailableStock
	p.LastUpdateDate = d.LastUpdateDate
	p.SupplierID = d.SupplierID
	p.ImageID = d.ImageID
}

type productsListResponse struct {
	Products []productDTO `json:"products"`
	Total    int          `json:"total"`
	Limit    int          `json:"limit"`
	Offset   int          `json:"offset"`
}
