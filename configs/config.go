package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config estructura de configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Keycloak KeycloakConfig
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig configuración de JWT
type JWTConfig struct {
	SecretKey     string
	AccessExpiry  int // en minutos
	RefreshExpiry int // en días
}

// KeycloakConfig configuración de Keycloak
type KeycloakConfig struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
	Enabled      bool
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Cargar archivo .env si existe
	godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "auth_service"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:     getEnv("JWT_SECRET_KEY", "your-secret-key"),
			AccessExpiry:  getEnvAsInt("JWT_ACCESS_EXPIRY", 15),
			RefreshExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY", 7),
		},
		Keycloak: KeycloakConfig{
			BaseURL:      getEnv("KEYCLOAK_BASE_URL", "http://localhost:8080"),
			Realm:        getEnv("KEYCLOAK_REALM", "master"),
			ClientID:     getEnv("KEYCLOAK_CLIENT_ID", "auth-service"),
			ClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", ""),
			Enabled:      getEnvAsBool("KEYCLOAK_ENABLED", false),
		},
	}

	return config, nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtiene una variable de entorno como entero o retorna un valor por defecto
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool obtiene una variable de entorno como booleano o retorna un valor por defecto
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetDSN retorna la cadena de conexión de la base de datos
func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}
