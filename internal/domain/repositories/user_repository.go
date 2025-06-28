package repositories

import (
	"context"

	"auth-go-microservicio/internal/domain/entities"
)

// UserRepository define las operaciones que debe implementar el repositorio de usuarios
type UserRepository interface {
	// Create crea un nuevo usuario
	Create(ctx context.Context, user *entities.User) error

	// GetByID obtiene un usuario por su ID
	GetByID(ctx context.Context, id string) (*entities.User, error)

	// GetByEmail obtiene un usuario por su email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *entities.User) error

	// Delete elimina un usuario por su ID
	Delete(ctx context.Context, id string) error

	// List obtiene una lista de usuarios con paginación
	List(ctx context.Context, offset, limit int) ([]*entities.User, error)

	// Count cuenta el total de usuarios
	Count(ctx context.Context) (int64, error)

	// ExistsByEmail verifica si existe un usuario con el email dado
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
