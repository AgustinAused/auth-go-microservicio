# 🔐 Microservicio de Autenticación y Autorización

Un microservicio robusto de autenticación y autorización construido en Go siguiendo los principios de **Clean Architecture**. Soporta autenticación local y **integración con Keycloak** como proveedor de identidad.

## ✨ Características

- 🔐 **Autenticación JWT** con tokens de acceso y refresh
- 🏗️ **Clean Architecture** con separación clara de responsabilidades
- 🛡️ **Autorización basada en roles** y permisos
- 🔄 **Integración con Keycloak** (Identity Provider)
- 📚 **Documentación Swagger** automática
- 🐳 **Docker y Docker Compose** listos para producción
- 🗄️ **PostgreSQL** como base de datos
- 🔒 **Hashing seguro** de contraseñas con bcrypt
- 🌐 **CORS** configurado
- 📝 **Logging** estructurado

## 🏗️ Arquitectura

### Sin Keycloak (Modo Local)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│  Auth Service   │───▶│   PostgreSQL    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   JWT Service   │
                       └─────────────────┘
```

### Con Keycloak (Modo Integrado)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│  Auth Service   │───▶│   Keycloak      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │  (Datos extra)  │
                       └─────────────────┘
```

## 🚀 Inicio Rápido

### Opción 1: Con Docker Compose (Recomendado)

```bash
# Clonar el repositorio
git clone <repository-url>
cd auth-go-microservicio

# Iniciar todos los servicios (incluyendo Keycloak)
docker-compose up -d

# Verificar servicios
docker-compose ps
```

**URLs disponibles:**
- 🔐 **Keycloak Admin Console**: http://localhost:8081
- 🌐 **API del Microservicio**: http://localhost:8080
- 📚 **Swagger Documentation**: http://localhost:8080/swagger/index.html

### Opción 2: Desarrollo Local

```bash
# Instalar dependencias
go mod tidy

# Configurar variables de entorno
cp env.example .env
# Editar .env según tus necesidades

# Ejecutar migraciones
make migrate

# Iniciar servidor
make run
```

## ⚙️ Configuración

### Variables de Entorno

```bash
# Configuración del servidor
SERVER_PORT=8080
SERVER_HOST=localhost

# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_service
DB_SSLMODE=disable

# JWT
JWT_SECRET_KEY=your-super-secret-jwt-key
JWT_ACCESS_EXPIRY=15
JWT_REFRESH_EXPIRY=7

# Keycloak (opcional)
KEYCLOAK_ENABLED=false
KEYCLOAK_BASE_URL=http://localhost:8081
KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=auth-service
KEYCLOAK_CLIENT_SECRET=your-keycloak-client-secret
```

## 🔐 Integración con Keycloak

### Habilitar Keycloak

1. **Configurar variables de entorno:**
```bash
KEYCLOAK_ENABLED=true
KEYCLOAK_BASE_URL=http://localhost:8081
KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=auth-service
KEYCLOAK_CLIENT_SECRET=your-secret
```

2. **Configurar Keycloak:**
   - Acceder a http://localhost:8081
   - Crear client `auth-service`
   - Configurar roles y usuarios
   - Ver [documentación completa](docs/KEYCLOAK_INTEGRATION.md)

### Autenticación con Keycloak

```bash
# Obtener token de Keycloak
curl -X POST http://localhost:8081/realms/master/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password&client_id=auth-service&client_secret=your-secret&username=admin&password=admin"

# Usar token en el microservicio
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📚 API Endpoints

### Autenticación (Públicos)
- `POST /api/v1/auth/register` - Registro de usuario
- `POST /api/v1/auth/login` - Login de usuario
- `POST /api/v1/auth/refresh` - Renovar token
- `POST /api/v1/auth/logout` - Logout de usuario

### Usuario (Requieren autenticación)
- `GET /api/v1/users/profile` - Obtener perfil
- `PUT /api/v1/users/profile` - Actualizar perfil
- `DELETE /api/v1/users/profile` - Eliminar cuenta
- `PUT /api/v1/users/change-password` - Cambiar contraseña

### Administración (Requieren rol admin)
- `GET /api/v1/admin/users` - Listar usuarios
- `PUT /api/v1/admin/users/{id}` - Actualizar usuario
- `DELETE /api/v1/admin/users/{id}` - Eliminar usuario

### Keycloak (Si está habilitado)
- `GET /api/v1/keycloak/users` - Listar usuarios de Keycloak
- `POST /api/v1/keycloak/users` - Crear usuario en Keycloak
- `PUT /api/v1/keycloak/users/{id}` - Actualizar usuario en Keycloak
- `DELETE /api/v1/keycloak/users/{id}` - Eliminar usuario de Keycloak
- `GET /api/v1/keycloak/users/{id}/groups` - Obtener grupos del usuario
- `PUT /api/v1/keycloak/users/{id}/groups/{group_id}` - Agregar usuario a grupo
- `DELETE /api/v1/keycloak/users/{id}/groups/{group_id}` - Remover usuario de grupo

## 🛠️ Comandos Útiles

```bash
# Ejecutar tests
make test

# Generar documentación Swagger
make swagger

# Ejecutar migraciones
make migrate

# Limpiar build
make clean

# Construir imagen Docker
make build

# Ejecutar con Docker
make docker-run
```

## 🏗️ Estructura del Proyecto

```
auth-go-microservicio/
├── cmd/
│   └── server/
│       └── main.go                 # Punto de entrada
├── configs/
│   └── config.go                   # Configuración
├── internal/
│   ├── domain/
│   │   ├── entities/               # Entidades del dominio
│   │   └── repositories/           # Interfaces de repositorios
│   ├── interface/
│   │   ├── database/
│   │   │   └── postgres/           # Implementación PostgreSQL
│   │   └── http/
│   │       ├── handlers/           # Manejadores HTTP
│   │       └── routes/             # Configuración de rutas
│   └── usecase/                    # Casos de uso
├── pkg/
│   ├── jwt/                        # Servicio JWT
│   ├── keycloak/                   # Servicio Keycloak
│   ├── middleware/                 # Middlewares
│   └── password/                   # Servicio de contraseñas
├── migrations/                     # Migraciones SQL
├── docs/                          # Documentación
├── docker-compose.yml             # Docker Compose
├── Dockerfile                     # Dockerfile
└── Makefile                       # Comandos útiles
```

## 🔧 Desarrollo

### Prerrequisitos
- Go 1.21+
- PostgreSQL 15+
- Docker y Docker Compose (opcional)
- Keycloak (opcional)

### Instalación Local

```bash
# Clonar repositorio
git clone <repository-url>
cd auth-go-microservicio

# Instalar dependencias
go mod download

# Configurar base de datos
# Crear base de datos PostgreSQL
createdb auth_service

# Ejecutar migraciones
make migrate

# Ejecutar tests
make test

# Iniciar servidor
make run
```

### Generar Documentación Swagger

```bash
# Instalar swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
make swagger
```

## 🧪 Testing

```bash
# Ejecutar todos los tests
make test

# Ejecutar tests con coverage
make test-coverage

# Ejecutar tests específicos
go test ./internal/usecase/...
```

## 🐳 Docker

### Construir Imagen

```bash
# Construir imagen
docker build -t auth-service .

# Ejecutar contenedor
docker run -p 8080:8080 auth-service
```

### Docker Compose

```bash
# Iniciar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```

## 📚 Documentación

- [📖 Guía de Integración con Keycloak](docs/KEYCLOAK_INTEGRATION.md)
- [📖 Configuración de Swagger](docs/SWAGGER_SETUP.md)
- [📖 API Documentation](docs/API.md)

## 🤝 Contribuir

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 🆘 Soporte

Si tienes problemas o preguntas:

1. Revisar la [documentación](docs/)
2. Buscar en [issues existentes](../../issues)
3. Crear un nuevo [issue](../../issues/new)

## 🔗 Enlaces Útiles

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [Gin Framework](https://gin-gonic.com/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [JWT.io](https://jwt.io/) 
