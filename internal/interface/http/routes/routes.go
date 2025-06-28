package routes

import (
	"auth-go-microservicio/configs"
	"auth-go-microservicio/internal/interface/http/handlers"
	"auth-go-microservicio/pkg/middleware"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	_ "auth-go-microservicio/docs" // Importar docs generados por swag

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	keycloakHandler *handlers.KeycloakHandler,
	authMiddleware *middleware.AuthMiddleware,
	keycloakMiddleware *middleware.KeycloakMiddleware,
	config *configs.Config,
) *gin.Engine {
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.Default())

	// Middleware de logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Rutas de autenticación (públicas)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/register-admin", authHandler.RegisterAdmin)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// Rutas de usuario (requieren autenticación)
		users := v1.Group("/users")
		users.Use(authMiddleware.Authenticate())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.DELETE("/profile", userHandler.DeleteAccount)
			users.PUT("/change-password", userHandler.ChangePassword)
		}

		// Rutas de administración (requieren rol de admin)
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.Authenticate())
		admin.Use(authMiddleware.RequireRole("admin"))
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.PUT("/users/:id", userHandler.UpdateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
		}

		// Rutas de Keycloak (si está habilitado)
		if config.Keycloak.Enabled {
			keycloak := v1.Group("/keycloak")
			keycloak.Use(keycloakMiddleware.Authenticate())
			keycloak.Use(keycloakMiddleware.RequireAdmin())
			{
				// Gestión de usuarios
				keycloak.GET("/users", keycloakHandler.GetUsers)
				keycloak.GET("/users/:id", keycloakHandler.GetUserByID)
				keycloak.POST("/users", keycloakHandler.CreateUser)
				keycloak.PUT("/users/:id", keycloakHandler.UpdateUser)
				keycloak.DELETE("/users/:id", keycloakHandler.DeleteUser)

				// Gestión de grupos
				keycloak.GET("/users/:id/groups", keycloakHandler.GetUserGroups)
				keycloak.PUT("/users/:id/groups/:group_id", keycloakHandler.AddUserToGroup)
				keycloak.DELETE("/users/:id/groups/:group_id", keycloakHandler.RemoveUserFromGroup)
			}
		}
	}

	// Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	return router
}
