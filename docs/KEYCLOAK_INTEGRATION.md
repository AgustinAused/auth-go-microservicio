# Integración con Keycloak

Este documento describe cómo integrar el microservicio de autenticación con Keycloak como proveedor de identidad (Identity Provider).

## 🏗️ Arquitectura

### Opción 1: Keycloak como Identity Provider (Recomendado)
- **Keycloak**: Maneja autenticación, autorización y gestión de usuarios
- **Microservicio**: Valida tokens JWT de Keycloak y proporciona endpoints adicionales
- **Base de datos**: Almacena datos adicionales del usuario (opcional)

### Opción 2: Híbrida
- **Keycloak**: Para autenticación y roles básicos
- **Microservicio**: Mantiene datos adicionales y lógica de negocio específica

## 🚀 Configuración

### 1. Variables de Entorno

```bash
# Configuración de Keycloak
KEYCLOAK_ENABLED=true
KEYCLOAK_BASE_URL=http://localhost:8081
KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=auth-service
KEYCLOAK_CLIENT_SECRET=your-keycloak-client-secret
```

### 2. Configuración de Keycloak

#### 2.1 Acceder a Keycloak Admin Console
- URL: `http://localhost:8081`
- Usuario: `admin`
- Contraseña: `admin`

#### 2.2 Crear un Realm (Opcional)
1. Ir a "Add realm"
2. Nombre: `auth-service`
3. Crear realm

#### 2.3 Crear un Client
1. Ir a "Clients" → "Create"
2. Client ID: `auth-service`
3. Client Protocol: `openid-connect`
4. Root URL: `http://localhost:8080`

#### 2.4 Configurar el Client
1. **Settings**:
   - Access Type: `confidential`
   - Valid Redirect URIs: `http://localhost:8080/*`
   - Web Origins: `http://localhost:8080`

2. **Credentials**:
   - Copiar el Client Secret y configurarlo en las variables de entorno

#### 2.5 Crear Roles
1. Ir a "Roles" → "Add Role"
   - `admin`: Para administradores
   - `user`: Para usuarios regulares
   - `moderator`: Para moderadores

#### 2.6 Crear Usuarios
1. Ir a "Users" → "Add User"
2. Configurar:
   - Username
   - Email
   - First Name
   - Last Name
3. En "Credentials":
   - Establecer contraseña
   - Desactivar "Temporary"
4. En "Role Mappings":
   - Asignar roles apropiados

## 🔧 Uso del Microservicio

### 1. Autenticación con Keycloak

#### Obtener Token de Acceso
```bash
curl -X POST http://localhost:8081/realms/master/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password&client_id=auth-service&client_secret=your-secret&username=admin&password=admin"
```

#### Usar Token en el Microservicio
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 2. Endpoints de Keycloak

#### Gestión de Usuarios
```bash
# Obtener todos los usuarios
GET /api/v1/keycloak/users

# Obtener usuario por ID
GET /api/v1/keycloak/users/{id}

# Crear usuario
POST /api/v1/keycloak/users
{
  "username": "newuser",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "enabled": true,
  "credentials": [
    {
      "type": "password",
      "value": "password123",
      "temporary": false
    }
  ]
}

# Actualizar usuario
PUT /api/v1/keycloak/users/{id}
{
  "firstName": "Jane",
  "lastName": "Smith"
}

# Eliminar usuario
DELETE /api/v1/keycloak/users/{id}
```

#### Gestión de Grupos
```bash
# Obtener grupos del usuario
GET /api/v1/keycloak/users/{id}/groups

# Agregar usuario a grupo
PUT /api/v1/keycloak/users/{id}/groups/{group_id}

# Remover usuario de grupo
DELETE /api/v1/keycloak/users/{id}/groups/{group_id}
```

## 🔐 Middleware de Autenticación

### Middleware de Keycloak
```go
// Autenticación básica
keycloakMiddleware.Authenticate()

// Verificar roles específicos
keycloakMiddleware.RequireRole("admin", "moderator")

// Verificar si es administrador
keycloakMiddleware.RequireAdmin()

// Verificar grupos
keycloakMiddleware.RequireGroup("developers", "admins")
```

### Información del Usuario en el Contexto
```go
// Obtener información del usuario autenticado
userID := c.GetString("user_id")
email := c.GetString("email")
username := c.GetString("username")
firstName := c.GetString("first_name")
lastName := c.GetString("last_name")
roles := c.GetStringSlice("realm_roles")
userInfo := c.MustGet("user_info").(*keycloak.UserInfo)
```

## 🐳 Ejecución con Docker Compose

### 1. Iniciar todos los servicios
```bash
docker-compose up -d
```

### 2. Verificar servicios
```bash
# Keycloak Admin Console
http://localhost:8081

# Microservicio API
http://localhost:8080

# Swagger Documentation
http://localhost:8080/swagger/index.html
```

### 3. Configurar Keycloak
1. Acceder a Keycloak Admin Console
2. Crear client y configurar según las instrucciones anteriores
3. Crear usuarios y roles
4. Actualizar variables de entorno con el client secret

## 🔄 Migración desde Autenticación Local

### 1. Habilitar Keycloak
```bash
KEYCLOAK_ENABLED=true
```

### 2. Migrar Usuarios
```bash
# Script de migración (ejemplo)
curl -X POST http://localhost:8080/api/v1/keycloak/users \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "existing_user",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "enabled": true,
    "credentials": [
      {
        "type": "password",
        "value": "migrated_password",
        "temporary": true
      }
    ]
  }'
```

### 3. Actualizar Frontend
```javascript
// Cambiar endpoint de autenticación
const loginUrl = 'http://localhost:8081/realms/master/protocol/openid-connect/token';

// Usar tokens de Keycloak
const token = response.access_token;
```

## 🛡️ Seguridad

### 1. Configuración de Seguridad
- Usar HTTPS en producción
- Configurar CORS apropiadamente
- Rotar client secrets regularmente
- Usar roles y grupos para autorización granular

### 2. Validación de Tokens
- Verificación de firma JWT
- Validación de expiración
- Verificación de issuer
- Validación de audience

### 3. Logs y Monitoreo
```bash
# Ver logs de Keycloak
docker-compose logs keycloak

# Ver logs del microservicio
docker-compose logs auth-service
```

## 🔧 Troubleshooting

### Problemas Comunes

#### 1. Error de Conexión a Keycloak
```bash
# Verificar que Keycloak esté ejecutándose
docker-compose ps keycloak

# Verificar logs
docker-compose logs keycloak
```

#### 2. Error de Validación de Token
```bash
# Verificar configuración del client
# Verificar client secret
# Verificar realm
```

#### 3. Error de CORS
```bash
# Configurar Web Origins en Keycloak
# Verificar configuración de CORS en el microservicio
```

### Logs de Depuración
```bash
# Habilitar logs detallados
export LOG_LEVEL=debug

# Ver logs en tiempo real
docker-compose logs -f auth-service
```

## 📚 Recursos Adicionales

- [Documentación oficial de Keycloak](https://www.keycloak.org/documentation)
- [OpenID Connect Specification](https://openid.net/connect/)
- [JWT.io](https://jwt.io/) - Para debuggear tokens JWT
- [Keycloak Admin REST API](https://www.keycloak.org/docs-api/24.0.2/rest-api/index.html) 