package services

import (
	"errors"

	"github.com/financial-analytics/api-gateway/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	ValidateToken(token string) (*Claims, error)
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type authService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) AuthService {
	return &authService{cfg: cfg}
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	// TODO: Implement actual token validation via Auth Service (gRPC/REST)
	// For now, this is a skeleton implementation to satisfy build requirements

	if tokenString == "valid-token" {
		return &Claims{
			UserID: "user-123",
			Email:  "user@example.com",
		}, nil
	}

	return nil, errors.New("invalid token")
}
