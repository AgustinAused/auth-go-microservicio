package handlers

import (
	"net/http"
	"strconv"

	"auth-go-microservicio/internal/usecase"

	"github.com/gin-gonic/gin"
)

// UserHandler maneja las peticiones HTTP de usuarios
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetProfile godoc
// @Summary      Obtener perfil del usuario
// @Description  Obtiene el perfil del usuario autenticado
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	req := &usecase.GetProfileRequest{
		UserID: userID.(string),
	}

	user, err := h.userUseCase.GetProfile(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile retrieved successfully",
		"data":    user,
	})
}

// UpdateProfile godoc
// @Summary      Actualizar perfil del usuario
// @Description  Actualiza el perfil del usuario autenticado
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body usecase.UpdateProfileRequest true "Datos del perfil"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req usecase.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID.(string)

	user, err := h.userUseCase.UpdateProfile(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"data":    user,
	})
}

// ChangePassword godoc
// @Summary      Cambiar contraseña
// @Description  Cambia la contraseña del usuario autenticado
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body usecase.ChangePasswordRequest true "Datos de cambio de contraseña"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/change-password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req usecase.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID.(string)

	err := h.userUseCase.ChangePassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password changed successfully",
	})
}

// DeleteAccount godoc
// @Summary      Eliminar cuenta
// @Description  Elimina la cuenta del usuario autenticado
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/profile [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	req := &usecase.DeleteAccountRequest{
		UserID: userID.(string),
	}

	err := h.userUseCase.DeleteAccount(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "account deleted successfully",
	})
}

// ListUsers godoc
// @Summary      Listar usuarios
// @Description  Lista todos los usuarios con paginación (solo para administradores)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        offset query int false "Offset para paginación" default(0)
// @Param        limit  query int false "Límite de resultados" default(10)
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &usecase.ListUsersRequest{
		Offset: offset,
		Limit:  limit,
	}

	response, err := h.userUseCase.ListUsers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "users retrieved successfully",
		"data":    response,
	})
}

// UpdateUser godoc
// @Summary      Actualizar usuario
// @Description  Actualiza un usuario específico (solo para administradores)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del usuario"
// @Param        request body usecase.UpdateUserRequest true "Datos del usuario"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	var req usecase.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID

	user, err := h.userUseCase.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"data":    user,
	})
}

// DeleteUser godoc
// @Summary      Eliminar usuario
// @Description  Elimina un usuario específico (solo para administradores)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID del usuario"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	req := &usecase.DeleteUserRequest{
		UserID: userID,
	}

	err := h.userUseCase.DeleteUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}
