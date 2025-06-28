package handlers

import (
	"net/http"

	"auth-go-microservicio/pkg/keycloak"

	"github.com/gin-gonic/gin"
)

// KeycloakHandler maneja las operaciones relacionadas con Keycloak
type KeycloakHandler struct {
	keycloakService keycloak.Service
}

// NewKeycloakHandler crea una nueva instancia del handler de Keycloak
func NewKeycloakHandler(keycloakService keycloak.Service) *KeycloakHandler {
	return &KeycloakHandler{
		keycloakService: keycloakService,
	}
}

// GetUsers godoc
// @Summary Obtener todos los usuarios de Keycloak
// @Description Obtiene la lista de todos los usuarios registrados en Keycloak
// @Tags keycloak
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} keycloak.UserInfo
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users [get]
func (h *KeycloakHandler) GetUsers(c *gin.Context) {
	users, err := h.keycloakService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting users from Keycloak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
		"count":   len(users),
	})
}

// GetUserByID godoc
// @Summary Obtener usuario por ID
// @Description Obtiene la información de un usuario específico de Keycloak por su ID
// @Tags keycloak
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Security BearerAuth
// @Success 200 {object} keycloak.UserInfo
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users/{id} [get]
func (h *KeycloakHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := h.keycloakService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user from Keycloak"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// CreateUser godoc
// @Summary Crear usuario en Keycloak
// @Description Crea un nuevo usuario en Keycloak
// @Tags keycloak
// @Accept json
// @Produce json
// @Param user body keycloak.CreateUserRequest true "Datos del usuario"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users [post]
func (h *KeycloakHandler) CreateUser(c *gin.Context) {
	var createUserReq keycloak.CreateUserRequest
	if err := c.ShouldBindJSON(&createUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validaciones básicas
	if createUserReq.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	if createUserReq.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	err := h.keycloakService.CreateUser(&createUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user in Keycloak"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "user created successfully",
	})
}

// UpdateUser godoc
// @Summary Actualizar usuario en Keycloak
// @Description Actualiza la información de un usuario existente en Keycloak
// @Tags keycloak
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param user body keycloak.UpdateUserRequest true "Datos actualizados del usuario"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users/{id} [put]
func (h *KeycloakHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var updateUserReq keycloak.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.keycloakService.UpdateUser(userID, &updateUserReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating user in Keycloak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user updated successfully",
	})
}

// DeleteUser godoc
// @Summary Eliminar usuario de Keycloak
// @Description Elimina un usuario de Keycloak
// @Tags keycloak
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users/{id} [delete]
func (h *KeycloakHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	err := h.keycloakService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting user from Keycloak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user deleted successfully",
	})
}

// GetUserGroups godoc
// @Summary Obtener grupos del usuario
// @Description Obtiene los grupos a los que pertenece un usuario en Keycloak
// @Tags keycloak
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Security BearerAuth
// @Success 200 {array} keycloak.Group
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users/{id}/groups [get]
func (h *KeycloakHandler) GetUserGroups(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	groups, err := h.keycloakService.GetUserGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user groups from Keycloak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    groups,
		"count":   len(groups),
	})
}

// AddUserToGroup godoc
// @Summary Agregar usuario a grupo
// @Description Agrega un usuario a un grupo específico en Keycloak
// @Tags keycloak
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param group_id path string true "ID del grupo"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users/{id}/groups/{group_id} [put]
func (h *KeycloakHandler) AddUserToGroup(c *gin.Context) {
	userID := c.Param("id")
	groupID := c.Param("group_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group ID is required"})
		return
	}

	err := h.keycloakService.AddUserToGroup(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error adding user to group in Keycloak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user added to group successfully",
	})
}

// RemoveUserFromGroup godoc
// @Summary Remover usuario de grupo
// @Description Remueve un usuario de un grupo específico en Keycloak
// @Tags keycloak
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Param group_id path string true "ID del grupo"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /keycloak/users/{id}/groups/{group_id} [delete]
func (h *KeycloakHandler) RemoveUserFromGroup(c *gin.Context) {
	userID := c.Param("id")
	groupID := c.Param("group_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group ID is required"})
		return
	}

	err := h.keycloakService.RemoveUserFromGroup(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error removing user from group in Keycloak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user removed from group successfully",
	})
}
