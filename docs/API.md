# API Documentation - Microservicio de Autenticación

## Base URL
```
http://localhost:8080/api/v1
```

## Endpoints

### Autenticación

#### 1. Registro de Usuario
**POST** `/auth/register`

Registra un nuevo usuario en el sistema.

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "contraseña123",
  "first_name": "Juan",
  "last_name": "Pérez"
}
```

**Response (201):**
```json
{
  "message": "user registered successfully",
  "data": {
    "user": {
      "id": "uuid-del-usuario",
      "email": "usuario@ejemplo.com",
      "first_name": "Juan",
      "last_name": "Pérez",
      "role": "user",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "token": "jwt-token-de-acceso"
  }
}
```

#### 2. Login de Usuario
**POST** `/auth/login`

Autentica un usuario y retorna tokens de acceso.

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "contraseña123"
}
```

**Response (200):**
```json
{
  "message": "login successful",
  "data": {
    "user": {
      "id": "uuid-del-usuario",
      "email": "usuario@ejemplo.com",
      "first_name": "Juan",
      "last_name": "Pérez",
      "role": "user",
      "is_active": true,
      "last_login_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "jwt-token-de-acceso",
    "refresh_token": "jwt-refresh-token"
  }
}
```

#### 3. Logout
**POST** `/auth/logout`

Cierra la sesión del usuario revocando el refresh token.

**Request Body:**
```json
{
  "refresh_token": "jwt-refresh-token"
}
```

**Response (200):**
```json
{
  "message": "logout successful"
}
```

#### 4. Refresh Token
**POST** `/auth/refresh`

Renueva el token de acceso usando un refresh token válido.

**Request Body:**
```json
{
  "refresh_token": "jwt-refresh-token"
}
```

**Response (200):**
```json
{
  "message": "token refreshed successfully",
  "data": {
    "access_token": "nuevo-jwt-token-de-acceso",
    "refresh_token": "nuevo-jwt-refresh-token"
  }
}
```

### Usuarios (Requiere Autenticación)

#### 1. Obtener Perfil
**GET** `/users/profile`

Obtiene el perfil del usuario autenticado.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200):**
```json
{
  "message": "profile retrieved successfully",
  "data": {
    "id": "uuid-del-usuario",
    "email": "usuario@ejemplo.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role": "user",
    "is_active": true,
    "last_login_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 2. Actualizar Perfil
**PUT** `/users/profile`

Actualiza el perfil del usuario autenticado.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "first_name": "Juan Carlos",
  "last_name": "Pérez García"
}
```

**Response (200):**
```json
{
  "message": "profile updated successfully",
  "data": {
    "id": "uuid-del-usuario",
    "email": "usuario@ejemplo.com",
    "first_name": "Juan Carlos",
    "last_name": "Pérez García",
    "role": "user",
    "is_active": true,
    "last_login_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3. Cambiar Contraseña
**PUT** `/users/change-password`

Cambia la contraseña del usuario autenticado.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "current_password": "contraseña123",
  "new_password": "nueva-contraseña456"
}
```

**Response (200):**
```json
{
  "message": "password changed successfully"
}
```

#### 4. Eliminar Cuenta
**DELETE** `/users/profile`

Elimina la cuenta del usuario autenticado.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200):**
```json
{
  "message": "account deleted successfully"
}
```

### Administración (Requiere Rol Admin)

#### 1. Listar Usuarios
**GET** `/admin/users?offset=0&limit=10`

Lista todos los usuarios con paginación.

**Headers:**
```
Authorization: Bearer <jwt-token-admin>
```

**Response (200):**
```json
{
  "message": "users retrieved successfully",
  "data": {
    "users": [
      {
        "id": "uuid-del-usuario",
        "email": "usuario@ejemplo.com",
        "first_name": "Juan",
        "last_name": "Pérez",
        "role": "user",
        "is_active": true,
        "last_login_at": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1
  }
}
```

#### 2. Actualizar Usuario
**PUT** `/admin/users/{id}`

Actualiza un usuario específico.

**Headers:**
```
Authorization: Bearer <jwt-token-admin>
```

**Request Body:**
```json
{
  "first_name": "Juan Carlos",
  "last_name": "Pérez García",
  "role": "admin",
  "is_active": true
}
```

**Response (200):**
```json
{
  "message": "user updated successfully",
  "data": {
    "id": "uuid-del-usuario",
    "email": "usuario@ejemplo.com",
    "first_name": "Juan Carlos",
    "last_name": "Pérez García",
    "role": "admin",
    "is_active": true,
    "last_login_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3. Eliminar Usuario
**DELETE** `/admin/users/{id}`

Elimina un usuario específico.

**Headers:**
```
Authorization: Bearer <jwt-token-admin>
```

**Response (200):**
```json
{
  "message": "user deleted successfully"
}
```

### Health Check

#### 1. Health Check
**GET** `/health`

Verifica el estado del servicio.

**Response (200):**
```json
{
  "status": "ok",
  "service": "auth-service"
}
```

## Códigos de Error

| Código | Descripción |
|--------|-------------|
| 400 | Bad Request - Datos de entrada inválidos |
| 401 | Unauthorized - Token inválido o faltante |
| 403 | Forbidden - Permisos insuficientes |
| 404 | Not Found - Recurso no encontrado |
| 500 | Internal Server Error - Error interno del servidor |

## Autenticación

El servicio utiliza JWT (JSON Web Tokens) para la autenticación. Los tokens deben incluirse en el header `Authorization` con el formato:

```
Authorization: Bearer <jwt-token>
```

## Roles

- **user**: Usuario normal con acceso a su propio perfil
- **admin**: Administrador con acceso a todos los usuarios

## Límites y Validaciones

- **Email**: Debe ser un email válido y único
- **Password**: Mínimo 8 caracteres
- **Nombres**: Máximo 100 caracteres cada uno
- **Paginación**: Offset y limit opcionales, por defecto limit=10 