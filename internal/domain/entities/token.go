package entities

import (
	"time"

	"github.com/google/uuid"
)

// Token representa un token JWT
type Token struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	TokenType TokenType `json:"token_type"`
	IsRevoked bool      `json:"is_revoked"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenType representa el tipo de token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// NewToken crea una nueva instancia de Token
func NewToken(userID uuid.UUID, token string, tokenType TokenType, expiresAt time.Time) *Token {
	return &Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		TokenType: tokenType,
		IsRevoked: false,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

// IsExpired verifica si el token ha expirado
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// Revoke marca el token como revocado
func (t *Token) Revoke() {
	t.IsRevoked = true
}

// IsValid verifica si el token es válido (no expirado y no revocado)
func (t *Token) IsValid() bool {
	return !t.IsExpired() && !t.IsRevoked
}
