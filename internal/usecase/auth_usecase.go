package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"auth-go-microservicio/internal/domain/entities"
	"auth-go-microservicio/internal/domain/repositories"
	"auth-go-microservicio/pkg/jwt"
	"auth-go-microservicio/pkg/keycloak"
	"auth-go-microservicio/pkg/password"
)

// AuthUseCase maneja la lógica de negocio para autenticación
type AuthUseCase struct {
	userRepo        repositories.UserRepository
	tokenRepo       repositories.TokenRepository
	jwtSvc          jwt.Service
	passSvc         password.Service
	keycloakService keycloak.Service
	keycloakConfig  *KeycloakConfig
	useKeycloak     bool
}

// KeycloakConfig configuración para Keycloak
type KeycloakConfig struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

// NewAuthUseCase crea una nueva instancia de AuthUseCase
func NewAuthUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	jwtSvc jwt.Service,
	passSvc password.Service,
	keycloakService keycloak.Service,
	keycloakConfig *KeycloakConfig,
) *AuthUseCase {
	// Determinar si usar Keycloak basado en la configuración
	useKeycloak := keycloakService != nil && keycloakConfig != nil &&
		keycloakConfig.BaseURL != "" && keycloakConfig.ClientID != "" && keycloakConfig.ClientSecret != ""

	return &AuthUseCase{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		jwtSvc:          jwtSvc,
		passSvc:         passSvc,
		keycloakService: keycloakService,
		keycloakConfig:  keycloakConfig,
		useKeycloak:     useKeycloak,
	}
}

// RegisterRequest representa la solicitud de registro
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"omitempty,oneof=user admin moderator"`
}

// RegisterResponse representa la respuesta del registro
type RegisterResponse struct {
	User  *entities.User `json:"user"`
	Token string         `json:"token"`
}

// Register registra un nuevo usuario
func (uc *AuthUseCase) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	if uc.useKeycloak {
		return uc.registerWithKeycloak(ctx, req)
	}
	return uc.registerLocal(ctx, req)
}

// registerWithKeycloak registra un usuario en Keycloak
func (uc *AuthUseCase) registerWithKeycloak(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Determinar el rol del usuario
	role := "user" // rol por defecto
	if req.Role != "" {
		switch req.Role {
		case "admin":
			role = "admin"
		case "moderator":
			role = "moderator"
		case "user":
			role = "user"
		default:
			return nil, errors.New("invalid role")
		}
	}

	// Crear usuario en Keycloak
	createUserReq := &keycloak.CreateUserRequest{
		Username:      req.Email,
		Email:         req.Email,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Enabled:       true,
		EmailVerified: true,
		Credentials: []*keycloak.Credential{
			{
				Type:      "password",
				Value:     req.Password,
				Temporary: false,
			},
		},
	}

	err := uc.keycloakService.CreateUser(createUserReq)
	if err != nil {
		return nil, fmt.Errorf("error creating user in Keycloak: %w", err)
	}

	// Obtener token de acceso para el usuario recién creado
	accessToken, err := uc.getKeycloakToken(req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("error getting access token: %w", err)
	}

	// Crear entidad de usuario local para respuesta
	user := &entities.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      entities.Role(role),
		IsActive:  true,
	}

	return &RegisterResponse{
		User:  user,
		Token: accessToken,
	}, nil
}

// registerLocal registra un usuario en la base de datos local
func (uc *AuthUseCase) registerLocal(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
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

	// Determinar el rol del usuario
	role := entities.RoleUser // rol por defecto
	if req.Role != "" {
		switch req.Role {
		case "admin":
			role = entities.RoleAdmin
		case "moderator":
			role = entities.RoleModerator
		case "user":
			role = entities.RoleUser
		default:
			return nil, errors.New("invalid role")
		}
	}

	// Crear usuario
	user := entities.NewUserWithRole(req.Email, hashedPassword, req.FirstName, req.LastName, role)

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
	if uc.useKeycloak {
		return uc.loginWithKeycloak(ctx, req)
	}
	return uc.loginLocal(ctx, req)
}

// loginWithKeycloak autentica un usuario usando Keycloak
func (uc *AuthUseCase) loginWithKeycloak(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Obtener token de acceso de Keycloak
	accessToken, err := uc.getKeycloakToken(req.Email, req.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Obtener información del usuario desde Keycloak
	userInfo, err := uc.keycloakService.GetUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("error getting user info: %w", err)
	}

	// Crear entidad de usuario local para respuesta
	user := &entities.User{
		Email:     userInfo.Email,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		IsActive:  userInfo.Enabled,
	}

	// Determinar rol del usuario (por defecto user)
	user.Role = entities.RoleUser
	if len(userInfo.Groups) > 0 {
		// Verificar si pertenece a grupos de admin
		for _, group := range userInfo.Groups {
			if group == "admin" || group == "administrators" {
				user.Role = entities.RoleAdmin
				break
			} else if group == "moderator" || group == "moderators" {
				user.Role = entities.RoleModerator
				break
			}
		}
	}

	// Para Keycloak, el refresh token se maneja automáticamente
	// No necesitamos almacenarlo localmente
	refreshToken := "" // Keycloak maneja esto internamente

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// loginLocal autentica un usuario usando la base de datos local
func (uc *AuthUseCase) loginLocal(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
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
	if uc.useKeycloak {
		// Para Keycloak, el logout se maneja automáticamente
		// No necesitamos hacer nada especial aquí
		return nil
	}
	// Revocar el refresh token local
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
	if uc.useKeycloak {
		return uc.refreshWithKeycloak(ctx, req)
	}
	return uc.refreshLocal(ctx, req)
}

// refreshWithKeycloak renueva un token usando Keycloak
func (uc *AuthUseCase) refreshWithKeycloak(ctx context.Context, req *RefreshRequest) (*RefreshResponse, error) {
	// Obtener nuevo token usando refresh token de Keycloak
	newAccessToken, err := uc.refreshKeycloakToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken, // Keycloak maneja el refresh token
	}, nil
}

// refreshLocal renueva un token usando la base de datos local
func (uc *AuthUseCase) refreshLocal(ctx context.Context, req *RefreshRequest) (*RefreshResponse, error) {
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

// getKeycloakToken obtiene un token de acceso de Keycloak
func (uc *AuthUseCase) getKeycloakToken(username, password string) (string, error) {
	return uc.keycloakService.Login(username, password)
}

// refreshKeycloakToken renueva un token usando el refresh token
func (uc *AuthUseCase) refreshKeycloakToken(refreshToken string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", uc.keycloakConfig.BaseURL, uc.keycloakConfig.Realm)

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", uc.keycloakConfig.ClientID)
	data.Set("client_secret", uc.keycloakConfig.ClientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error refreshing token: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// IsUsingKeycloak retorna si el sistema está usando Keycloak
func (uc *AuthUseCase) IsUsingKeycloak() bool {
	return uc.useKeycloak
}
