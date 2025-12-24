package jwt

import (
	"time"

	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateToken(userID uuid.UUID, role Role) (string, error)
}

type TokenValidator interface {
	ValidateToken(tokenString string) (*Claims, error)
}

type Manager struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

func NewJWTManager(secret string, ttl time.Duration, issuer string) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
		issuer: issuer,
	}
}

func (m *Manager) GenerateToken(userID uuid.UUID, role Role) (string, error) {
	return GenerateToken(userID, role, m.secret, m.ttl, m.issuer)
}

func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	return ParseToken(tokenString, m.secret, m.issuer)
}
