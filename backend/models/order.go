package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID              uint           `json:"id" gorm:"primaryKey" example:"1"`
	OrderNumber     string         `json:"order_number" gorm:"unique;not null" example:"BJJ-1753519000"`
	GuestEmail      string         `json:"guest_email" gorm:"not null" example:"customer@example.com"`
	ShippingAddress Address        `json:"shipping_address" gorm:"type:jsonb"`
	Items           []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount     float64        `json:"total_amount" gorm:"not null" example:"120.00"`
	Status          OrderStatus    `json:"status" gorm:"default:pending" example:"pending"`
	StripePaymentID string         `json:"stripe_payment_id" example:"pi_1234567890"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	OrderID   uint      `json:"order_id" gorm:"not null" example:"1"`
	ProductID uint      `json:"product_id" gorm:"not null" example:"1"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int       `json:"quantity" gorm:"not null" example:"1"`
	Price     float64   `json:"price" gorm:"not null" example:"120.00"`
	Size      string    `json:"size" example:"A2"`
	CreatedAt time.Time `json:"created_at"`
}

type Address struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Address1  string `json:"address1" example:"123 Main St"`
	Address2  string `json:"address2" example:"Apt 4B"`
	City      string `json:"city" example:"Ho Chi Minh City"`
	State     string `json:"state" example:"HCMC"`
	ZipCode   string `json:"zip_code" example:"70000"`
	Country   string `json:"country" example:"Vietnam"`
}

// Add Cart Item for frontend
type CartItem struct {
	ProductID uint    `json:"product_id" example:"1"`
	Quantity  int     `json:"quantity" example:"1"`
	Size      string  `json:"size" example:"A2"`
	Price     float64 `json:"price" example:"120.00"`
}

// GORM JSON handling for Address
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Address) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &a)
}

// Business methods
func (o *Order) CalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.TotalAmount = total
}

func (o *Order) GenerateOrderNumber() string {
	return fmt.Sprintf("BJJ-%d", time.Now().Unix())
}
