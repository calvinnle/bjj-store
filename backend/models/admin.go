package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminRole string

const (
	RoleSuperAdmin   AdminRole = "super_admin"
	RoleInventory    AdminRole = "inventory"
	RoleOrderManager AdminRole = "order_manager"
	RoleViewer       AdminRole = "viewer"
)

type AdminUser struct {
	ID           uint           `json:"id" gorm:"primaryKey" example:"1"`
	Email        string         `json:"email" gorm:"unique;not null" example:"admin@bjjstore.com"`
	PasswordHash string         `json:"-" gorm:"not null"` // Don't include in JSON
	Role         AdminRole      `json:"role" gorm:"default:viewer" example:"super_admin"`
	IsActive     bool           `json:"is_active" gorm:"default:true" example:"true"`
	LastLogin    *time.Time     `json:"last_login" example:"2025-07-26T15:20:41.306413+07:00"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type AdminSession struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	AdminID   uint      `json:"admin_id" gorm:"not null" example:"1"`
	TokenHash string    `json:"-" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false" example:"false"`
	CreatedAt time.Time `json:"created_at"`
}

// Hash password before saving
func (a *AdminUser) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.PasswordHash = string(hash)
	return nil
}

// Check if password is correct
func (a *AdminUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	return err == nil
}

// Check if admin has permission
func (a *AdminUser) HasPermission(permission string) bool {
	switch a.Role {
	case RoleSuperAdmin:
		return true // Super admin has all permissions
	case RoleInventory:
		return permission == "view_products" || permission == "create_products" ||
			permission == "update_products" || permission == "delete_products"
	case RoleOrderManager:
		return permission == "view_orders" || permission == "update_orders"
	case RoleViewer:
		return permission == "view_products" || permission == "view_orders"
	default:
		return false
	}
}
