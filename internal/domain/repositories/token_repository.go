package repositories

import (
	"context"

	"auth-go-microservicio/internal/domain/entities"
)

// TokenRepository define las operaciones que debe implementar el repositorio de tokens
type TokenRepository interface {
	// Create crea un nuevo token
	Create(ctx context.Context, token *entities.Token) error

	// GetByToken obtiene un token por su valor
	GetByToken(ctx context.Context, token string) (*entities.Token, error)

	// GetByUserID obtiene todos los tokens de un usuario
	GetByUserID(ctx context.Context, userID string) ([]*entities.Token, error)

	// RevokeByUserID revoca todos los tokens de un usuario
	RevokeByUserID(ctx context.Context, userID string) error

	// RevokeToken revoca un token específico
	RevokeToken(ctx context.Context, token string) error

	// DeleteExpired elimina tokens expirados
	DeleteExpired(ctx context.Context) error

	// Cleanup limpia tokens antiguos y revocados
	Cleanup(ctx context.Context) error
}
