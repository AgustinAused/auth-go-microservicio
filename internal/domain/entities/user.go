package entities

import (
	"time"

	"github.com/google/uuid"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	Password    string     `json:"-"` // No se serializa en JSON
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Role        Role       `json:"role"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Role representa los roles de usuario
type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
)

// NewUser crea una nueva instancia de User
func NewUser(email, password, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      RoleUser, // Por defecto es usuario
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewUserWithRole crea una nueva instancia de User con un rol específico
func NewUserWithRole(email, password, firstName, lastName string, role Role) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// FullName retorna el nombre completo del usuario
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin verifica si el usuario es administrador
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// UpdateLastLogin actualiza la fecha del último login
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// Deactivate desactiva el usuario
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activa el usuario
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}
