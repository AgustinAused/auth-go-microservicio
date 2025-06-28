# Variables
BINARY_NAME=auth-service
BUILD_DIR=build
DOCKER_IMAGE=auth-service
DOCKER_TAG=latest

# Comandos principales
.PHONY: help build run test clean docker-build docker-run docker-stop

help: ## Muestra esta ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compila la aplicación
	@echo "Compilando la aplicación..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

run: ## Ejecuta la aplicación localmente
	@echo "Ejecutando la aplicación..."
	go run ./cmd/server/main.go

test: ## Ejecuta los tests
	@echo "Ejecutando tests..."
	go test -v ./...

test-coverage: ## Ejecuta los tests con cobertura
	@echo "Ejecutando tests con cobertura..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Limpia archivos generados
	@echo "Limpiando archivos..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

deps: ## Instala dependencias
	@echo "Instalando dependencias..."
	go mod tidy
	go mod download

# Swagger commands
swagger-init: ## Genera la documentación de Swagger
	@echo "Generando documentación de Swagger..."
	swag init -g cmd/server/main.go

swagger-serve: ## Sirve la documentación de Swagger
	@echo "Sirviendo documentación de Swagger..."
	swag serve -F=swagger docs/swagger.json

# Docker commands
docker-build: ## Construye la imagen Docker
	@echo "Construyendo imagen Docker..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Ejecuta la aplicación con Docker Compose
	@echo "Ejecutando con Docker Compose..."
	docker-compose up -d

docker-stop: ## Detiene la aplicación con Docker Compose
	@echo "Deteniendo Docker Compose..."
	docker-compose down

docker-logs: ## Muestra logs de Docker Compose
	docker-compose logs -f

# Database commands
db-migrate: ## Ejecuta migraciones de la base de datos
	@echo "Ejecutando migraciones..."
	@echo "Asegúrate de tener PostgreSQL ejecutándose y configurado en .env"

db-reset: ## Resetea la base de datos (¡CUIDADO!)
	@echo "¡ADVERTENCIA! Esto eliminará todos los datos."
	@read -p "¿Estás seguro? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo "Reseteando base de datos..."

# Development commands
dev-setup: ## Configura el entorno de desarrollo
	@echo "Configurando entorno de desarrollo..."
	@if [ ! -f .env ]; then cp env.example .env; echo "Archivo .env creado desde env.example"; fi
	@echo "Asegúrate de configurar las variables en .env"

lint: ## Ejecuta el linter
	@echo "Ejecutando linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint no está instalado. Instalando..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

format: ## Formatea el código
	@echo "Formateando código..."
	go fmt ./...

# Production commands
prod-build: ## Construye para producción
	@echo "Construyendo para producción..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Utility commands
check-env: ## Verifica que el archivo .env existe
	@if [ ! -f .env ]; then \
		echo "Error: Archivo .env no encontrado"; \
		echo "Ejecuta 'make dev-setup' para crear el archivo"; \
		exit 1; \
	fi

health-check: ## Verifica el estado del servicio
	@echo "Verificando estado del servicio..."
	@curl -f http://localhost:8080/health || echo "Servicio no está ejecutándose"

swagger-check: ## Verifica que Swagger UI esté disponible
	@echo "Verificando Swagger UI..."
	@curl -f http://localhost:8080/swagger/index.html || echo "Swagger UI no está disponible"

# Default target
.DEFAULT_GOAL := help 