// @title           Microservicio de Autenticación API
// @version         1.0
// @description     Un microservicio robusto de autenticación y autorización construido en Go siguiendo los principios de Clean Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"auth-go-microservicio/configs"
	"auth-go-microservicio/internal/domain/repositories"
	"auth-go-microservicio/internal/interface/database/postgres"
	"auth-go-microservicio/internal/interface/http/handlers"
	"auth-go-microservicio/internal/interface/http/routes"
	"auth-go-microservicio/internal/usecase"
	"auth-go-microservicio/pkg/jwt"
	"auth-go-microservicio/pkg/keycloak"
	"auth-go-microservicio/pkg/middleware"
	"auth-go-microservicio/pkg/password"

	_ "github.com/lib/pq"

	_ "auth-go-microservicio/docs" // Importar docs generados por swag
)

// @title           Microservicio de Autenticación API
// @version         1.0
// @description     Un microservicio robusto de autenticación y autorización construido en Go siguiendo los principios de Clean Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Cargar configuración
	config, err := configs.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Conectar a la base de datos
	db, err := sql.Open("postgres", config.GetDSN())
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Verificar conexión a la base de datos
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	// Inicializar servicios
	jwtService := jwt.NewService(
		config.JWT.SecretKey,
		time.Duration(config.JWT.AccessExpiry)*time.Minute,
		time.Duration(config.JWT.RefreshExpiry)*24*time.Hour,
	)
	passwordService := password.NewService(12) // bcrypt cost 12

	// Inicializar repositorios
	var userRepo repositories.UserRepository
	var tokenRepo repositories.TokenRepository

	// Si Keycloak está habilitado, usar repositorios de Keycloak
	if config.Keycloak.Enabled {
		keycloakService := keycloak.NewService(
			config.Keycloak.BaseURL,
			config.Keycloak.Realm,
			config.Keycloak.ClientID,
			config.Keycloak.ClientSecret,
		)

		// Para Keycloak, podríamos usar repositorios híbridos o solo Keycloak
		// Por ahora, mantenemos los repositorios de PostgreSQL para datos adicionales
		userRepo = postgres.NewUserRepository(db)
		tokenRepo = postgres.NewTokenRepository(db)

		// Inicializar use cases
		authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, passwordService)
		userUseCase := usecase.NewUserUseCase(userRepo, passwordService)

		// Inicializar middlewares
		authMiddleware := middleware.NewAuthMiddleware(jwtService)
		keycloakMiddleware := middleware.NewKeycloakMiddleware(keycloakService)

		// Inicializar handlers
		authHandler := handlers.NewAuthHandler(authUseCase)
		userHandler := handlers.NewUserHandler(userUseCase)
		keycloakHandler := handlers.NewKeycloakHandler(keycloakService)

		// Configurar rutas
		router := routes.SetupRoutes(authHandler, userHandler, keycloakHandler, authMiddleware, keycloakMiddleware, config)

		// Iniciar servidor
		serverAddr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
		log.Printf("🚀 Servidor iniciado en %s", serverAddr)
		log.Printf("📚 Swagger UI disponible en http://%s/swagger/index.html", serverAddr)
		log.Printf("🔐 Keycloak habilitado - Realm: %s", config.Keycloak.Realm)

		if err := http.ListenAndServe(serverAddr, router); err != nil {
			log.Fatal("Error starting server:", err)
		}
	} else {
		// Modo sin Keycloak (funcionalidad original)
		userRepo = postgres.NewUserRepository(db)
		tokenRepo = postgres.NewTokenRepository(db)

		// Inicializar use cases
		authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, passwordService)
		userUseCase := usecase.NewUserUseCase(userRepo, passwordService)

		// Inicializar middleware
		authMiddleware := middleware.NewAuthMiddleware(jwtService)

		// Inicializar handlers
		authHandler := handlers.NewAuthHandler(authUseCase)
		userHandler := handlers.NewUserHandler(userUseCase)

		// Configurar rutas (sin Keycloak)
		router := routes.SetupRoutes(authHandler, userHandler, nil, authMiddleware, nil, config)

		// Iniciar servidor
		serverAddr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
		log.Printf("🚀 Servidor iniciado en %s", serverAddr)
		log.Printf("📚 Swagger UI disponible en http://%s/swagger/index.html", serverAddr)
		log.Printf("🔐 Modo autenticación local (sin Keycloak)")

		if err := http.ListenAndServe(serverAddr, router); err != nil {
			log.Fatal("Error starting server:", err)
		}
	}
}
