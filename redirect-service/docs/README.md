# API Documentation

This directory contains auto-generated Swagger/OpenAPI documentation for the Re:Path Redirect Service API.

## 📚 Swagger Files

- **`swagger.json`** - OpenAPI specification in JSON format
- **`swagger.yaml`** - OpenAPI specification in YAML format
- **`docs.go`** - Go package with embedded documentation

## 🚀 Accessing the Documentation

### Swagger UI (Interactive)

When the service is running, access the interactive Swagger UI at:

```
http://localhost:8081/swagger/index.html
```

The Swagger UI provides:
- Interactive API testing
- Request/response examples
- Schema definitions
- Try it out functionality

### Alternative Ports

If you've configured a different port via `APP_PORT` environment variable:

```
http://localhost:{YOUR_PORT}/swagger/index.html
```

## 🔄 Regenerating Documentation

After modifying API endpoints or adding new handlers with Swagger annotations, regenerate the docs:

```bash
# Using Make
make swagger

# Or directly
swag init -g cmd/main.go -o docs
```

## 📝 Adding Swagger Annotations

### Handler Example

```go
// GetURLInfo returns URL information without redirecting
// @Summary Get URL information
// @Description Get detailed information about a shortened URL
// @Tags Redirect
// @Accept json
// @Produce json
// @Param shortCode path string true "Short Code"
// @Success 200 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/info/{shortCode} [get]
func (h *RedirectHandler) GetURLInfo(c *gin.Context) {
    // Handler implementation
}
```

### Main API Info

The main API information is defined in `cmd/main.go`:

```go
// @title           Re:Path Redirect Service API
// @version         1.0
// @description     A high-performance URL shortening service
// @host            localhost:8081
// @BasePath        /
```

## 📖 Available Endpoints

### Health Check
- **GET** `/health` - Service health check

### Redirect Operations
- **GET** `/{shortCode}` - Redirect to original URL (301/302)
- **GET** `/api/info/{shortCode}` - Get URL information without redirecting

## 🔗 Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [gin-swagger](https://github.com/swaggo/gin-swagger)

## 🛠️ Development

### Auto-regenerate on changes

The docs are automatically regenerated when you:
1. Modify handler annotations
2. Run `make swagger`
3. Rebuild the application

### Best Practices

1. **Always document new endpoints** - Add Swagger annotations to all new handlers
2. **Keep descriptions clear** - Write concise, helpful descriptions
3. **Define response models** - Use proper DTO structs for responses
4. **Test in Swagger UI** - Use the UI to test your endpoints
5. **Regenerate after changes** - Always run `make swagger` after modifying APIs

## 🎯 Quick Start

```bash
# 1. Install swag tool
make install

# 2. Generate docs
make swagger

# 3. Start the service
make dev

# 4. Open Swagger UI
open http://localhost:8081/swagger/index.html
```

---

**Note:** These files are auto-generated. Do not edit them directly. Instead, modify the Swagger annotations in your Go code and regenerate.
