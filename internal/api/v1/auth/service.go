package auth

import (
	"tyr/internal/repo"

	"github.com/M15t/gram/pkg/server/middleware/jwt"

	gjwt "github.com/golang-jwt/jwt/v5"
)

// New creates new auth service
func New(repo *repo.Service, jwt JWT, cr Crypter) *Auth {
	return &Auth{
		repo: repo,
		jwt:  jwt,
		cr:   cr,
	}
}

// Auth represents auth application service
type Auth struct {
	repo *repo.Service
	jwt  JWT
	cr   Crypter
}

// JWT represents token generator (jwt) interface
type JWT interface {
	GenerateToken(*jwt.TokenInput, *jwt.TokenOutput) error
	ParseToken(string) (*gjwt.Token, error)
}

// Crypter represents security interface
type Crypter interface {
	HashPassword(string) string
	CompareHashAndPassword(string, string) bool
}
