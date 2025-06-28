package usecase

import (
	"context"
	"errors"

	"auth-go-microservicio/internal/domain/entities"
	"auth-go-microservicio/internal/domain/repositories"
	"auth-go-microservicio/pkg/password"
)

// UserUseCase maneja la lógica de negocio para usuarios
type UserUseCase struct {
	userRepo repositories.UserRepository
	passSvc  password.Service
}

// NewUserUseCase crea una nueva instancia de UserUseCase
func NewUserUseCase(userRepo repositories.UserRepository, passSvc password.Service) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		passSvc:  passSvc,
	}
}

// GetProfileRequest representa la solicitud para obtener perfil
type GetProfileRequest struct {
	UserID string
}

// GetProfile obtiene el perfil de un usuario
func (uc *UserUseCase) GetProfile(ctx context.Context, req *GetProfileRequest) (*entities.User, error) {
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateProfileRequest representa la solicitud para actualizar perfil
type UpdateProfileRequest struct {
	UserID    string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateProfile actualiza el perfil de un usuario
func (uc *UserUseCase) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*entities.User, error) {
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Actualizar campos
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	// Guardar cambios
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePasswordRequest representa la solicitud para cambiar contraseña
type ChangePasswordRequest struct {
	UserID          string `json:"-"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword cambia la contraseña de un usuario
func (uc *UserUseCase) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verificar contraseña actual
	if !uc.passSvc.Verify(req.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash de la nueva contraseña
	hashedPassword, err := uc.passSvc.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	// Actualizar contraseña
	user.Password = hashedPassword
	return uc.userRepo.Update(ctx, user)
}

// DeleteAccountRequest representa la solicitud para eliminar cuenta
type DeleteAccountRequest struct {
	UserID string `json:"-"`
}

// DeleteAccount elimina la cuenta de un usuario
func (uc *UserUseCase) DeleteAccount(ctx context.Context, req *DeleteAccountRequest) error {
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.userRepo.Delete(ctx, user.ID.String())
}

// ListUsersRequest representa la solicitud para listar usuarios
type ListUsersRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ListUsersResponse representa la respuesta para listar usuarios
type ListUsersResponse struct {
	Users []*entities.User `json:"users"`
	Total int64            `json:"total"`
}

// ListUsers lista usuarios (solo para administradores)
func (uc *UserUseCase) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	users, err := uc.userRepo.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &ListUsersResponse{
		Users: users,
		Total: total,
	}, nil
}

// UpdateUserRequest representa la solicitud para actualizar usuario (admin)
type UpdateUserRequest struct {
	UserID    string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	IsActive  *bool  `json:"is_active"`
}

// UpdateUser actualiza un usuario (solo para administradores)
func (uc *UserUseCase) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*entities.User, error) {
	user, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Actualizar campos
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Role != "" {
		user.Role = entities.Role(req.Role)
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Guardar cambios
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUserRequest representa la solicitud para eliminar usuario (admin)
type DeleteUserRequest struct {
	UserID string `json:"-"`
}

// DeleteUser elimina un usuario (solo para administradores)
func (uc *UserUseCase) DeleteUser(ctx context.Context, req *DeleteUserRequest) error {
	return uc.userRepo.Delete(ctx, req.UserID)
}
