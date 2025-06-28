package handlers

import (
	"net/http"

	"auth-go-microservicio/internal/usecase"

	"github.com/gin-gonic/gin"
)

// AuthHandler maneja las peticiones HTTP de autenticación
type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

// NewAuthHandler crea una nueva instancia de AuthHandler
func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Register godoc
// @Summary      Registro de usuario
// @Description  Registra un nuevo usuario en el sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body usecase.RegisterRequest true "Datos de registro"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"data":    response,
	})
}

// Login godoc
// @Summary      Login de usuario
// @Description  Autentica un usuario y retorna tokens de acceso
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body usecase.LoginRequest true "Credenciales de login"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data":    response,
	})
}

// Logout godoc
// @Summary      Logout de usuario
// @Description  Cierra la sesión del usuario revocando el refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body usecase.LogoutRequest true "Refresh token"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req usecase.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authUseCase.Logout(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

// Refresh godoc
// @Summary      Refresh token
// @Description  Renueva el token de acceso usando un refresh token válido
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body usecase.RefreshRequest true "Refresh token"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req usecase.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authUseCase.Refresh(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed successfully",
		"data":    response,
	})
}
