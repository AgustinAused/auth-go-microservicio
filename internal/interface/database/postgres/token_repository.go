package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"auth-go-microservicio/internal/domain/entities"
	"auth-go-microservicio/internal/domain/repositories"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// TokenRepository implementa el repositorio de tokens para PostgreSQL
type TokenRepository struct {
	db *sql.DB
}

// NewTokenRepository crea una nueva instancia de TokenRepository
func NewTokenRepository(db *sql.DB) repositories.TokenRepository {
	return &TokenRepository{db: db}
}

// Create crea un nuevo token
func (r *TokenRepository) Create(ctx context.Context, token *entities.Token) error {
	query := `
		INSERT INTO tokens (id, user_id, token, token_type, is_revoked, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.TokenType,
		token.IsRevoked,
		token.ExpiresAt,
		token.CreatedAt,
	)

	return err
}

// GetByToken obtiene un token por su valor
func (r *TokenRepository) GetByToken(ctx context.Context, token string) (*entities.Token, error) {
	query := `
		SELECT id, user_id, token, token_type, is_revoked, expires_at, created_at
		FROM tokens WHERE token = $1
	`

	var tokenEntity entities.Token

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&tokenEntity.ID,
		&tokenEntity.UserID,
		&tokenEntity.Token,
		&tokenEntity.TokenType,
		&tokenEntity.IsRevoked,
		&tokenEntity.ExpiresAt,
		&tokenEntity.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	return &tokenEntity, nil
}

// GetByUserID obtiene todos los tokens de un usuario
func (r *TokenRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.Token, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	query := `
		SELECT id, user_id, token, token_type, is_revoked, expires_at, created_at
		FROM tokens WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, parsedUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entities.Token

	for rows.Next() {
		var token entities.Token

		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.TokenType,
			&token.IsRevoked,
			&token.ExpiresAt,
			&token.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		tokens = append(tokens, &token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

// RevokeByUserID revoca todos los tokens de un usuario
func (r *TokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	query := `UPDATE tokens SET is_revoked = true WHERE user_id = $1`

	_, err = r.db.ExecContext(ctx, query, parsedUserID)
	return err
}

// RevokeToken revoca un token específico
func (r *TokenRepository) RevokeToken(ctx context.Context, token string) error {
	query := `UPDATE tokens SET is_revoked = true WHERE token = $1`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}

// DeleteExpired elimina tokens expirados
func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM tokens WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	return err
}

// Cleanup limpia tokens antiguos y revocados
func (r *TokenRepository) Cleanup(ctx context.Context) error {
	// Eliminar tokens revocados más antiguos de 30 días
	query := `
		DELETE FROM tokens 
		WHERE is_revoked = true AND created_at < $1
	`

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	_, err := r.db.ExecContext(ctx, query, thirtyDaysAgo)

	return err
}
