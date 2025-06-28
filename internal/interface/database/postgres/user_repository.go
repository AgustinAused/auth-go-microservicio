package postgres

import (
	"context"
	"database/sql"
	"errors"

	"auth-go-microservicio/internal/domain/entities"
	"auth-go-microservicio/internal/domain/repositories"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// UserRepository implementa el repositorio de usuarios para PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &UserRepository{db: db}
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	query := `
		SELECT id, email, password, first_name, last_name, role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user entities.User
	var lastLoginAt sql.NullTime

	err = r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, is_active, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user entities.User
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $2, password = $3, first_name = $4, last_name = $5, role = $6, 
		    is_active = $7, last_login_at = $8, updated_at = $9
		WHERE id = $1
	`

	var lastLoginAt sql.NullTime
	if user.LastLoginAt != nil {
		lastLoginAt.Time = *user.LastLoginAt
		lastLoginAt.Valid = true
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		lastLoginAt,
		user.UpdatedAt,
	)

	return err
}

// Delete elimina un usuario por su ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user id")
	}

	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// List obtiene una lista de usuarios con paginación
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, is_active, last_login_at, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User

	for rows.Next() {
		var user entities.User
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.IsActive,
			&lastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Count cuenta el total de usuarios
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)

	return count, err
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)

	return exists, err
}
