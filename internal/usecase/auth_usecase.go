package usecase

import (
	"context"
	"errors"
	"time"

	"auth-go-microservicio/internal/domain/entities"
	"auth-go-microservicio/internal/domain/repositories"
	"auth-go-microservicio/pkg/jwt"
	"auth-go-microservicio/pkg/password"
)

// AuthUseCase maneja la lógica de negocio para autenticación
type AuthUseCase struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.TokenRepository
	jwtSvc    jwt.Service
	passSvc   password.Service
}

// NewAuthUseCase crea una nueva instancia de AuthUseCase
func NewAuthUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	jwtSvc jwt.Service,
	passSvc password.Service,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSvc:    jwtSvc,
		passSvc:   passSvc,
	}
}

// RegisterRequest representa la solicitud de registro
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// RegisterResponse representa la respuesta del registro
type RegisterResponse struct {
	User  *entities.User `json:"user"`
	Token string         `json:"token"`
}

// Register registra un nuevo usuario
func (uc *AuthUseCase) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Verificar si el email ya existe
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash de la contraseña
	hashedPassword, err := uc.passSvc.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// Crear usuario
	user := entities.NewUser(req.Email, hashedPassword, req.FirstName, req.LastName)

	// Guardar en la base de datos
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generar token JWT
	token, err := uc.jwtSvc.GenerateToken(user.ID.String(), user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		User:  user,
		Token: token,
	}, nil
}

// LoginRequest representa la solicitud de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse representa la respuesta del login
type LoginResponse struct {
	User         *entities.User `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
}

// Login autentica un usuario
func (uc *AuthUseCase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Obtener usuario por email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Verificar contraseña
	if !uc.passSvc.Verify(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Actualizar último login
	user.UpdateLastLogin()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Generar tokens
	accessToken, err := uc.jwtSvc.GenerateToken(user.ID.String(), user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtSvc.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	// Guardar refresh token en la base de datos
	refreshTokenEntity := entities.NewToken(
		user.ID,
		refreshToken,
		entities.TokenTypeRefresh,
		time.Now().Add(24*7*time.Hour), // 7 días
	)
	if err := uc.tokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LogoutRequest representa la solicitud de logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Logout cierra la sesión del usuario
func (uc *AuthUseCase) Logout(ctx context.Context, req *LogoutRequest) error {
	// Revocar el refresh token
	return uc.tokenRepo.RevokeToken(ctx, req.RefreshToken)
}

// RefreshRequest representa la solicitud de refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse representa la respuesta del refresh
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Refresh renueva el token de acceso
func (uc *AuthUseCase) Refresh(ctx context.Context, req *RefreshRequest) (*RefreshResponse, error) {
	// Verificar el refresh token
	claims, err := uc.jwtSvc.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Verificar si el token existe en la base de datos
	token, err := uc.tokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Verificar si el token es válido
	if !token.IsValid() {
		return nil, errors.New("invalid refresh token")
	}

	// Obtener usuario
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Generar nuevos tokens
	accessToken, err := uc.jwtSvc.GenerateToken(user.ID.String(), user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwtSvc.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	// Revocar el token anterior
	if err := uc.tokenRepo.RevokeToken(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	// Guardar el nuevo refresh token
	newRefreshTokenEntity := entities.NewToken(
		user.ID,
		newRefreshToken,
		entities.TokenTypeRefresh,
		time.Now().Add(24*7*time.Hour), // 7 días
	)
	if err := uc.tokenRepo.Create(ctx, newRefreshTokenEntity); err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
