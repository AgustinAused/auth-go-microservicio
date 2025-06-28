package routes

import (
	"auth-go-microservicio/internal/interface/http/handlers"
	"auth-go-microservicio/pkg/middleware"

	_ "auth-go-microservicio/docs" // Importar docs generados por swag

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
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
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// Rutas de usuarios (requieren autenticación)
		users := v1.Group("/users")
		users.Use(authMiddleware.Authenticate())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.DELETE("/profile", userHandler.DeleteAccount)
		}

		// Rutas de administración (requieren rol admin)
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.Authenticate())
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.PUT("/users/:id", userHandler.UpdateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
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
