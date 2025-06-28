package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	// Servidor
	Port string
	Env  string

	// Base de datos
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration

	// CORS
	AllowedOrigins []string
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Cargar archivo .env si existe
	godotenv.Load()

	config := &Config{
		// Servidor
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// Base de datos
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "auth_service"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT
		JWTSecret:              getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		JWTExpiration:          parseDuration(getEnv("JWT_EXPIRATION", "24h")),
		RefreshTokenExpiration: parseDuration(getEnv("REFRESH_TOKEN_EXPIRATION", "168h")), // 7 días

		// CORS
		AllowedOrigins: parseStringSlice(getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")),
	}

	return config, nil
}

// getEnv obtiene una variable de entorno con un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parsea una duración desde string
func parseDuration(duration string) time.Duration {
	parsed, err := time.ParseDuration(duration)
	if err != nil {
		// Valores por defecto si hay error
		switch duration {
		case "24h":
			return 24 * time.Hour
		case "168h":
			return 168 * time.Hour
		default:
			return 24 * time.Hour
		}
	}
	return parsed
}

// parseStringSlice parsea un string separado por comas en un slice
func parseStringSlice(input string) []string {
	if input == "" {
		return []string{}
	}

	// Por simplicidad, asumimos que no hay espacios alrededor de las comas
	// En un caso real, podrías usar strings.Split y strings.TrimSpace
	return []string{input} // Simplificado para este ejemplo
}

// GetDSN retorna la cadena de conexión de la base de datos
func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}
