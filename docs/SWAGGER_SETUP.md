# Configuración y Uso de Swagger

## ¿Qué es Swagger?

Swagger (OpenAPI) es una herramienta que permite documentar automáticamente las APIs REST. En este proyecto, hemos integrado Swagger para generar documentación interactiva de nuestra API de autenticación.

## Características Implementadas

- ✅ Documentación automática de endpoints
- ✅ Interfaz web interactiva
- ✅ Autenticación JWT integrada
- ✅ Ejemplos de requests y responses
- ✅ Códigos de error documentados
- ✅ Parámetros y tipos de datos definidos

## URLs de Acceso

Una vez que el servidor esté ejecutándose:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Base**: http://localhost:8080/api/v1
- **Health Check**: http://localhost:8080/health

## Comandos Disponibles

### Generar Documentación
```bash
# Generar documentación desde los comentarios
make swagger-init
# o
swag init -g cmd/server/main.go
```

### Verificar Swagger
```bash
# Verificar que Swagger UI esté disponible
make swagger-check
```

### Servir Documentación Independiente
```bash
# Servir solo la documentación (sin el servidor principal)
make swagger-serve
```

## Estructura de Comentarios

### Comentarios Principales (main.go)
```go
// @title           Microservicio de Autenticación API
// @version         1.0
// @description     Descripción de la API
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Comentarios de Endpoints
```go
// @Summary      Título del endpoint
// @Description  Descripción detallada
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RequestType true "Descripción del parámetro"
// @Success      200  {object}  ResponseType
// @Failure      400  {object}  ErrorResponse
// @Router       /endpoint [method]
```

## Endpoints Documentados

### Autenticación
- `POST /auth/register` - Registro de usuario
- `POST /auth/login` - Login de usuario
- `POST /auth/logout` - Logout de usuario
- `POST /auth/refresh` - Refresh token

### Usuarios (Requiere Autenticación)
- `GET /users/profile` - Obtener perfil
- `PUT /users/profile` - Actualizar perfil
- `DELETE /users/profile` - Eliminar cuenta

### Administración (Requiere Rol Admin)
- `GET /admin/users` - Listar usuarios
- `PUT /admin/users/{id}` - Actualizar usuario
- `DELETE /admin/users/{id}` - Eliminar usuario

## Autenticación en Swagger

1. **Obtener Token**: Usa el endpoint `/auth/login` para obtener un token JWT
2. **Autorizar**: Haz clic en el botón "Authorize" en la parte superior de Swagger UI
3. **Ingresar Token**: En el campo "Value", ingresa: `Bearer <tu-token-jwt>`
4. **Probar Endpoints**: Ahora puedes probar endpoints protegidos

## Personalización

### Cambiar Información de la API
Edita los comentarios en `cmd/server/main.go`:

```go
// @title           Tu Título de API
// @version         2.0
// @description     Tu descripción personalizada
// @contact.name    Tu Nombre
// @contact.email   tu@email.com
```

### Agregar Nuevos Endpoints
1. Agrega comentarios de Swagger al handler
2. Ejecuta `make swagger-init` para regenerar la documentación
3. Reinicia el servidor

### Cambiar Tema de Swagger UI
Puedes personalizar el tema modificando las opciones en `routes.go`:

```go
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
    ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
    ginSwagger.DefaultModelsExpandDepth(-1),
))
```

## Troubleshooting

### Error: "swag command not found"
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Error: "docs package not found"
```bash
make swagger-init
```

### Error: "Swagger UI not loading"
1. Verifica que el servidor esté ejecutándose
2. Verifica que la ruta `/swagger/*any` esté configurada
3. Revisa los logs del servidor

### Documentación no se actualiza
1. Ejecuta `make swagger-init`
2. Reinicia el servidor
3. Limpia la caché del navegador

## Archivos Generados

- `docs/docs.go` - Código Go generado
- `docs/swagger.json` - Especificación JSON
- `docs/swagger.yaml` - Especificación YAML

## Mejores Prácticas

1. **Mantén la documentación actualizada**: Ejecuta `make swagger-init` después de cambios
2. **Usa descripciones claras**: Los comentarios deben ser descriptivos
3. **Documenta errores**: Incluye todos los códigos de error posibles
4. **Usa tags**: Agrupa endpoints relacionados con tags
5. **Ejemplos reales**: Proporciona ejemplos de requests y responses

## Recursos Adicionales

- [Documentación oficial de swaggo](https://github.com/swaggo/swag)
- [Especificación OpenAPI](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/) 