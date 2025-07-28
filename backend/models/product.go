package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	Name        string         `json:"name" gorm:"not null" example:"Tatami Estilo 6.0 Gi"`
	Description string         `json:"description" example:"Premium BJJ gi with excellent fit and durability"`
	Price       float64        `json:"price" gorm:"not null" example:"120.00"`
	Category    string         `json:"category" example:"gi"`
	SizeOptions string         `json:"size_options" gorm:"type:text" example:"A1,A2,A3,A4"` // Changed to simple string
	Stock       int            `json:"stock" gorm:"default:0" example:"15"`
	ImageURL    string         `json:"image_url" example:"https://example.com/gi.jpg"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Business methods
func (p *Product) IsAvailable() bool {
	return p.Stock > 0
}

func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return fmt.Errorf("insufficient stock: requested %d, available %d", quantity, p.Stock)
	}
	p.Stock -= quantity
	return nil
}

func (p *Product) HasSufficientStock(quantity int) bool {
	return p.Stock >= quantity
}

// Helper methods to work with size options as comma-separated string
func (p *Product) GetSizeOptionsArray() []string {
	if p.SizeOptions == "" {
		return []string{}
	}
	// Split by comma and trim spaces
	sizes := []string{}
	for _, size := range strings.Split(p.SizeOptions, ",") {
		sizes = append(sizes, strings.TrimSpace(size))
	}
	return sizes
}

func (p *Product) SetSizeOptionsFromArray(sizes []string) {
	p.SizeOptions = strings.Join(sizes, ",")
}
