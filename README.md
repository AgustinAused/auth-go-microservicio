# Microservicio de Autenticación y Autorización

Un microservicio robusto de autenticación y autorización construido en Go siguiendo los principios de Clean Architecture.

## Características

- ✅ Autenticación con JWT
- ✅ Registro de usuarios
- ✅ Login/Logout
- ✅ Autorización basada en roles
- ✅ Refresh tokens
- ✅ Validación de contraseñas seguras
- ✅ Middleware de autenticación
- ✅ Clean Architecture
- ✅ Base de datos PostgreSQL/MySQL
- ✅ Variables de entorno
- ✅ CORS habilitado
- ✅ Logging estructurado

## Estructura del Proyecto

```
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   └── repositories/
│   ├── usecase/
│   ├── interface/
│   │   ├── http/
│   │   └── database/
│   └── infrastructure/
├── pkg/
│   ├── jwt/
│   ├── password/
│   └── middleware/
├── configs/
├── migrations/
└── docs/
```

## Instalación

1. Clona el repositorio
2. Instala las dependencias:
   ```bash
   go mod tidy
   ```
3. Configura las variables de entorno:
   ```bash
   cp .env.example .env
   ```
4. Ejecuta las migraciones de la base de datos
5. Inicia el servidor:
   ```bash
   go run cmd/server/main.go
   ```

## Endpoints

### Autenticación
- `POST /api/v1/auth/register` - Registro de usuario
- `POST /api/v1/auth/login` - Login de usuario
- `POST /api/v1/auth/logout` - Logout de usuario
- `POST /api/v1/auth/refresh` - Refresh token

### Usuarios
- `GET /api/v1/users/profile` - Obtener perfil del usuario
- `PUT /api/v1/users/profile` - Actualizar perfil del usuario
- `DELETE /api/v1/users/profile` - Eliminar cuenta

### Admin (solo para administradores)
- `GET /api/v1/admin/users` - Listar todos los usuarios
- `PUT /api/v1/admin/users/:id` - Actualizar usuario
- `DELETE /api/v1/admin/users/:id` - Eliminar usuario

## Variables de Entorno

```env
# Servidor
PORT=8080
ENV=development

# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_service
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION=24h
REFRESH_TOKEN_EXPIRATION=168h

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

## Tecnologías Utilizadas

- **Go 1.21+** - Lenguaje principal
- **Gin** - Framework web
- **JWT** - Tokens de autenticación
- **PostgreSQL/MySQL** - Base de datos
- **bcrypt** - Hash de contraseñas
- **UUID** - Identificadores únicos

## Arquitectura

Este proyecto sigue los principios de Clean Architecture:

1. **Domain Layer**: Contiene las entidades y reglas de negocio
2. **Use Case Layer**: Implementa los casos de uso de la aplicación
3. **Interface Layer**: Maneja la comunicación externa (HTTP, DB)
4. **Infrastructure Layer**: Implementaciones concretas de las interfaces

## Contribución

1. Fork el proyecto
2. Crea una rama para tu feature
3. Commit tus cambios
4. Push a la rama
5. Abre un Pull Request

## Licencia

MIT 