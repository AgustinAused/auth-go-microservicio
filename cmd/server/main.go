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
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-go-microservicio/configs"
	"auth-go-microservicio/internal/interface/database/postgres"
	"auth-go-microservicio/internal/interface/http/handlers"
	"auth-go-microservicio/internal/interface/http/routes"
	"auth-go-microservicio/internal/usecase"
	"auth-go-microservicio/pkg/jwt"
	"auth-go-microservicio/pkg/middleware"
	"auth-go-microservicio/pkg/password"

	_ "github.com/lib/pq"
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
	jwtService := jwt.NewService(config.JWTSecret, config.JWTExpiration, config.RefreshTokenExpiration)
	passwordService := password.NewService(12) // bcrypt cost 12

	// Inicializar repositorios
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	// Inicializar casos de uso
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, passwordService)
	userUseCase := usecase.NewUserUseCase(userRepo, passwordService)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)

	// Inicializar middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Configurar rutas
	router := routes.SetupRoutes(authHandler, userHandler, authMiddleware)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para manejar señales de terminación
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("Starting server on port %s", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error starting server:", err)
		}
	}()

	// Esperar señal de terminación
	<-done
	log.Println("Shutting down server...")

	// Crear contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Intentar cerrar el servidor gracefulmente
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
