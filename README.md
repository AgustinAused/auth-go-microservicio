# 🔐 Microservicio de Autenticación y Autorización

Un microservicio robusto de autenticación y autorización construido en Go siguiendo los principios de **Clean Architecture**. Soporta **dos modos de autenticación**: local (JWT) y Keycloak como Identity Provider.

## 🚀 Características

### ✅ Funcionalidades Principales
- **Autenticación dual**: Local (JWT) o Keycloak
- **Autorización por roles**: admin, moderator, user
- **Gestión de usuarios**: registro, login, logout, refresh tokens
- **Middleware de autenticación**: flexible y configurable
- **Documentación automática**: Swagger/OpenAPI
- **Base de datos**: PostgreSQL con migraciones
- **Docker**: Contenedores listos para producción

### 🔄 Modos de Autenticación

#### 1. **Modo Local (JWT)**
- Usuarios almacenados en PostgreSQL
- Tokens JWT generados localmente
- Gestión completa de refresh tokens
- Ideal para aplicaciones simples

#### 2. **Modo Keycloak**
- Usuarios gestionados por Keycloak
- Tokens JWT generados por Keycloak
- Integración completa con Identity Provider
- Ideal para aplicaciones empresariales

## 🏗️ Arquitectura

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Cliente       │    │  Microservicio   │    │   Keycloak      │
│                 │    │                  │    │   (Opcional)    │
│ Login/Register  │───▶│ AuthUseCase      │◄──▶│                 │
│ (email/pass)    │    │ (Dual Mode)      │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │  PostgreSQL      │
                       │  (Datos locales) │
                       └──────────────────┘
```

## 🛠️ Instalación

### Prerrequisitos
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (opcional)
- Keycloak (solo si usas modo Keycloak)

### 1. Clonar el repositorio
```bash
git clone https://github.com/tu-usuario/auth-go-microservicio.git
cd auth-go-microservicio
```

### 2. Configurar variables de entorno
```bash
cp env.example .env
```

#### Para modo local (JWT):
```bash
# Configuración básica
SERVER_PORT=8080
DB_HOST=localhost
DB_PASSWORD=password
JWT_SECRET_KEY=your-super-secret-jwt-key

# Deshabilitar Keycloak
KEYCLOAK_ENABLED=false
```

#### Para modo Keycloak:
```bash
# Configuración básica
SERVER_PORT=8080
DB_HOST=localhost
DB_PASSWORD=password

# Habilitar Keycloak
KEYCLOAK_ENABLED=true
KEYCLOAK_BASE_URL=http://localhost:8081
KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=auth-service
KEYCLOAK_CLIENT_SECRET=your-secret
```

### 3. Ejecutar con Docker Compose
```bash
# Incluye PostgreSQL y Keycloak (opcional)
docker-compose up -d
```

### 4. Ejecutar localmente
```bash
# Instalar dependencias
go mod download

# Ejecutar migraciones
make migrate

# Iniciar servidor
make run
```

## 📚 API Endpoints

### Autenticación (Públicos)
- `POST /api/v1/auth/register` - Registro de usuario
- `POST /api/v1/auth/register-admin` - Registro de administrador
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

### Keycloak (Solo si está habilitado)
- `GET /api/v1/keycloak/users` - Listar usuarios de Keycloak
- `POST /api/v1/keycloak/users` - Crear usuario en Keycloak
- `PUT /api/v1/keycloak/users/{id}` - Actualizar usuario en Keycloak
- `DELETE /api/v1/keycloak/users/{id}` - Eliminar usuario de Keycloak

## 🔧 Configuración de Keycloak

Si usas el modo Keycloak, sigue estos pasos:

### 1. Acceder a Keycloak Admin Console
- URL: `http://localhost:8081`
- Usuario: `admin`
- Contraseña: `admin`

### 2. Crear Client
1. Ir a "Clients" → "Create"
2. Client ID: `auth-service`
3. Client Protocol: `openid-connect`
4. Root URL: `http://localhost:8080`

### 3. Configurar Client
- Access Type: `confidential`
- Valid Redirect URIs: `http://localhost:8080/*`
- Web Origins: `http://localhost:8080`

### 4. Obtener Client Secret
- Ir a "Credentials"
- Copiar el Client Secret
- Configurarlo en `KEYCLOAK_CLIENT_SECRET`

Ver [documentación completa](docs/KEYCLOAK_INTEGRATION.md) para más detalles.

## 🧪 Pruebas

### Ejecutar script de pruebas
```bash
# PowerShell
.\scripts\test-endpoints.ps1

# Bash
./scripts/test-endpoints.sh
```

### Probar manualmente
```bash
# 1. Registrar usuario
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 3. Usar token
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 📖 Documentación

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Health Check**: `http://localhost:8080/health`
- **Keycloak Admin**: `http://localhost:8081` (solo modo Keycloak)

## 🛠️ Comandos Útiles

```bash
# Desarrollo
make run          # Ejecutar servidor
make build        # Compilar
make test         # Ejecutar tests
make migrate      # Ejecutar migraciones

# Docker
make docker-build # Construir imagen
make docker-run   # Ejecutar con Docker

# Documentación
make swagger      # Generar documentación Swagger
```

## 🔄 Migración entre Modos

### De Local a Keycloak
1. Configurar Keycloak según la documentación
2. Establecer `KEYCLOAK_ENABLED=true`
3. Configurar variables de Keycloak
4. Reiniciar el servidor

### De Keycloak a Local
1. Establecer `KEYCLOAK_ENABLED=false`
2. Configurar `JWT_SECRET_KEY`
3. Reiniciar el servidor

El sistema detecta automáticamente qué modo usar basado en la configuración.

## 🤝 Contribuir

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request


## 🆘 Soporte

- 📧 Email: support@example.com
- 📖 Documentación: [docs/](docs/)
- 🐛 Issues: [GitHub Issues](https://github.com/tu-usuario/auth-go-microservicio/issues)
