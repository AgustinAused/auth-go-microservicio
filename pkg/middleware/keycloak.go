package middleware

import (
	"net/http"
	"strings"

	"auth-go-microservicio/pkg/keycloak"

	"github.com/gin-gonic/gin"
)

// KeycloakMiddleware middleware para autenticación con Keycloak
type KeycloakMiddleware struct {
	keycloakService keycloak.Service
}

// NewKeycloakMiddleware crea una nueva instancia del middleware de Keycloak
func NewKeycloakMiddleware(keycloakService keycloak.Service) *KeycloakMiddleware {
	return &KeycloakMiddleware{
		keycloakService: keycloakService,
	}
}

// Authenticate middleware para verificar autenticación con Keycloak
func (m *KeycloakMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validar token con Keycloak
		claims, err := m.keycloakService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Obtener información adicional del usuario
		userInfo, err := m.keycloakService.GetUserInfo(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error getting user info"})
			c.Abort()
			return
		}

		// Agregar información al contexto
		c.Set("user_id", claims.Sub)
		c.Set("email", claims.Email)
		c.Set("username", claims.PreferredUsername)
		c.Set("first_name", claims.GivenName)
		c.Set("last_name", claims.FamilyName)
		c.Set("realm_roles", claims.RealmAccess.Roles)
		c.Set("user_info", userInfo)
		c.Set("keycloak_claims", claims)

		c.Next()
	}
}

// RequireRole middleware para verificar roles específicos de Keycloak
func (m *KeycloakMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmRoles, exists := c.Get("realm_roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user roles not found"})
			c.Abort()
			return
		}

		userRoles, ok := realmRoles.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid roles type"})
			c.Abort()
			return
		}

		// Verificar si el usuario tiene alguno de los roles requeridos
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware para verificar que el usuario sea administrador
func (m *KeycloakMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin", "realm-admin")
}

// RequireGroup middleware para verificar que el usuario pertenezca a un grupo específico
func (m *KeycloakMiddleware) RequireGroup(groups ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, exists := c.Get("user_info")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user info not found"})
			c.Abort()
			return
		}

		info, ok := userInfo.(*keycloak.UserInfo)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user info type"})
			c.Abort()
			return
		}

		// Verificar si el usuario pertenece a alguno de los grupos requeridos
		hasGroup := false
		for _, requiredGroup := range groups {
			for _, userGroup := range info.Groups {
				if userGroup == requiredGroup {
					hasGroup = true
					break
				}
			}
			if hasGroup {
				break
			}
		}

		if !hasGroup {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient group permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
