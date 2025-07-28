package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/calvinnle/bjj-store/backend/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	AdminID uint             `json:"admin_id"`
	Email   string           `json:"email"`
	Role    models.AdminRole `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey string
}

// Remove the config dependency from constructor
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
	}
}

func (s *JWTService) GenerateToken(admin *models.AdminUser) (string, error) {
	claims := JWTClaims{
		AdminID: admin.ID,
		Email:   admin.Email,
		Role:    admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 2 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bjj-store",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *JWTService) StoreSession(adminID uint, token string) error {
	tokenHash := s.HashToken(token)

	session := models.AdminSession{
		AdminID:   adminID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(2 * time.Hour),
		IsRevoked: false,
	}

	return models.DB.Create(&session).Error
}

func (s *JWTService) IsTokenRevoked(token string) bool {
	tokenHash := s.HashToken(token)

	var session models.AdminSession
	err := models.DB.Where("token_hash = ? AND is_revoked = false AND expires_at > ?",
		tokenHash, time.Now()).First(&session).Error

	return err != nil // If not found, consider it revoked
}

func (s *JWTService) RevokeToken(token string) error {
	tokenHash := s.HashToken(token)

	return models.DB.Model(&models.AdminSession{}).
		Where("token_hash = ?", tokenHash).
		Update("is_revoked", true).Error
}
