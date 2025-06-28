package middleware

import (
	"net/http"
	"strings"

	"auth-go-microservicio/pkg/jwt"
	"auth-go-microservicio/pkg/keycloak"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware para autenticación
type AuthMiddleware struct {
	jwtService      jwt.Service
	keycloakService keycloak.Service
	useKeycloak     bool
}

// NewAuthMiddleware crea una nueva instancia del middleware de autenticación
func NewAuthMiddleware(jwtService jwt.Service, keycloakService keycloak.Service, useKeycloak bool) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:      jwtService,
		keycloakService: keycloakService,
		useKeycloak:     useKeycloak,
	}
}

// Authenticate middleware para verificar autenticación
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>" o solo el token
		var token string
		parts := strings.Split(authHeader, " ")

		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		} else if len(parts) == 1 {
			// Si solo se proporciona el token sin "Bearer"
			token = parts[0]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format. Use 'Bearer <token>' or just the token"})
			c.Abort()
			return
		}

		if m.useKeycloak && m.keycloakService != nil {
			// Usar autenticación de Keycloak
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

			// Determinar rol del usuario
			role := "user" // por defecto
			if len(claims.RealmAccess.Roles) > 0 {
				for _, r := range claims.RealmAccess.Roles {
					if r == "admin" || r == "realm-admin" {
						role = "admin"
						break
					} else if r == "moderator" {
						role = "moderator"
						break
					}
				}
			}
			c.Set("role", role)

		} else {
			// Usar autenticación JWT local
			claims, err := m.jwtService.ValidateToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}

			// Agregar claims al contexto
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}

// RequireRole middleware para verificar roles específicos
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid role type"})
			c.Abort()
			return
		}

		// Verificar si el rol del usuario está en la lista de roles permitidos
		hasRole := false
		for _, role := range roles {
			if roleStr == role {
				hasRole = true
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
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}
