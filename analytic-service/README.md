# Analytic Service

**Re:Path Analytics Service** - FastAPI microservice for analytics built with Python and FastAPI.

## Features

- **FastAPI Framework** - Modern, fast web framework for building APIs
- **Global Response Format** - Standardized API responses across all endpoints
- **Health Check Endpoint** - Kubernetes/Docker ready health monitoring
- **Pydantic V2** - Data validation using Python type annotations
- **Environment-based Configuration** - Secure configuration management
- **CORS Support** - Cross-Origin Resource Sharing enabled
- **Auto-generated API Docs** - Interactive Swagger UI and ReDoc
- **Docker Support** - Production-ready containerization
- **Exception Handling** - Centralized error handling with detailed error responses
- **Best Practices** - Clean architecture, type hints, async/await

## Project Structure

```
analytic-service/
├── app/
│   ├── __init__.py
│   ├── main.py              # FastAPI application entry point
│   ├── api/
│   │   └── v1/
│   │       ├── __init__.py
│   │       └── health.py    # Health check endpoint
│   ├── core/
│   │   ├── __init__.py
│   │   ├── config.py        # Configuration management
│   │   └── exceptions.py    # Custom exceptions & handlers
│   ├── models/              # Database models (future)
│   ├── schemas/             # Pydantic schemas
│   │   ├── __init__.py
│   │   ├── health.py        # Health check schemas
│   │   └── response.py      # Global response schemas
│   └── services/            # Business logic (future)
├── tests/                   # Test files
├── .env.example             # Example environment variables
├── .gitignore
├── .dockerignore
├── Dockerfile               # Multi-stage production build
├── docker-compose.yml       # Container orchestration
├── Makefile                 # Development commands
├── pyproject.toml           # Project configuration
├── requirements.txt         # Production dependencies
├── requirements-dev.txt     # Development dependencies
└── README.md
```

## Prerequisites

- Python 3.9 or higher
- pip (Python package installer)
- Docker & Docker Compose (optional, for containerized deployment)

> **Note**: The codebase is compatible with Python 3.9+ and uses `typing.Union` and `typing.List` for type hints to ensure backward compatibility.

## Quick Start

### 1. Local Development

```bash
# Install dependencies
make install

# Copy environment file
cp .env.example .env

# Run development server with hot-reload
make dev
```

The service will be available at:
- **API**: http://localhost:8000
- **Interactive API Docs**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc

### 2. Docker Development

```bash
# Build and run with Docker Compose
make docker-run

# Or run with hot-reload for development
make docker-dev

# View logs
make docker-logs

# Stop containers
make docker-stop
```

## Available Commands

### Development Commands

```bash
make help          # Show all available commands
make install       # Install production dependencies
make install-dev   # Install development dependencies
make dev           # Run development server with hot-reload
make run           # Run production server
make test          # Run tests with coverage
make lint          # Run linting checks
make format        # Format code with ruff
make type-check    # Run type checking with mypy
make clean         # Clean temporary files and caches
```

### Docker Commands

```bash
make docker-build  # Build Docker image
make docker-run    # Run Docker container (production)
make docker-dev    # Run Docker container with hot-reload
make docker-stop   # Stop Docker containers
make docker-logs   # View container logs
```

## API Endpoints

All API responses follow a standardized format. See [RESPONSE_EXAMPLES.md](./RESPONSE_EXAMPLES.md) for detailed examples.

### Health Check

**GET** `/api/v1/health`

Returns service health status and information in a standardized response format.

**Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "service": "analytic-service",
    "version": "0.1.0",
    "timestamp": "2025-10-14T10:30:00Z",
    "environment": "development"
  },
  "meta": {
    "timestamp": "2025-10-14T10:30:00Z",
    "version": "v1",
    "request_id": null
  }
}
```

### Global Response Format

All endpoints return responses in a consistent format:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* your data */ },
  "meta": {
    "timestamp": "2025-10-14T10:30:00Z",
    "version": "v1",
    "request_id": "req_abc123xyz"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error message",
  "errors": [
    {
      "code": "ERROR_CODE",
      "message": "Detailed error message",
      "field": "fieldName"
    }
  ],
  "meta": {
    "timestamp": "2025-10-14T10:30:00Z",
    "version": "v1",
    "request_id": "req_abc123xyz"
  }
}
```

See [RESPONSE_EXAMPLES.md](./RESPONSE_EXAMPLES.md) for more detailed examples including pagination, validation errors, and client usage.

## Configuration

Configuration is managed through environment variables. Copy `.env.example` to `.env` and adjust as needed:

```env
# Application Settings
APP_NAME=analytic-service
APP_VERSION=0.1.0
APP_ENV=development

# Server Configuration
HOST=0.0.0.0
PORT=8000
RELOAD=true

# API Configuration
API_V1_PREFIX=/api/v1

# CORS
CORS_ORIGINS=["http://localhost:3000","http://localhost:8000"]
CORS_ALLOW_CREDENTIALS=true
CORS_ALLOW_METHODS=["*"]
CORS_ALLOW_HEADERS=["*"]

# Logging
LOG_LEVEL=info
```

## Testing

```bash
# Run tests with coverage
make test

# Run specific test file
source venv/bin/activate
pytest tests/test_health.py -v
```

## Code Quality

```bash
# Run linting
make lint

# Format code
make format

# Type checking
make type-check

# Run all checks
make lint && make type-check && make test
```

## Best Practices Implemented

### Application Architecture
- **Clean Architecture**: Separation of concerns with distinct layers (API, schemas, services, models)
- **Dependency Injection**: Centralized configuration management
- **Type Safety**: Full type hints with mypy validation
- **Async/Await**: Asynchronous request handling for better performance

### Docker Best Practices
- **Multi-stage Build**: Smaller production images
- **Non-root User**: Enhanced security
- **Layer Caching**: Optimized build times
- **Health Checks**: Container health monitoring
- **.dockerignore**: Minimal image size

### Code Quality
- **Ruff**: Fast Python linter and formatter
- **MyPy**: Static type checking
- **Pytest**: Comprehensive testing framework
- **Code Coverage**: Track test coverage metrics

### API Design
- **RESTful**: Following REST principles
- **Versioning**: API versioning support (v1)
- **Documentation**: Auto-generated OpenAPI/Swagger docs
- **Validation**: Request/response validation with Pydantic

### Security
- **Environment Variables**: Secure configuration management
- **CORS**: Configurable CORS policies
- **Non-root Docker User**: Principle of least privilege
- **Input Validation**: Pydantic schema validation

## Development Workflow

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature
   ```

2. **Make Changes**
   - Write code following project structure
   - Add type hints to all functions
   - Write tests for new features

3. **Run Quality Checks**
   ```bash
   make format    # Format code
   make lint      # Check for issues
   make type-check # Verify types
   make test      # Run tests
   ```

4. **Test Locally**
   ```bash
   make dev       # Test with hot-reload
   ```

5. **Commit and Push**
   ```bash
   git add .
   git commit -m "feat: your feature description"
   git push origin feature/your-feature
   ```

## Production Deployment

### Using Docker

```bash
# Build production image
docker build -t analytic-service:latest .

# Run production container
docker run -d \
  -p 8000:8000 \
  --env-file .env \
  --name analytic-service \
  analytic-service:latest
```

### Using Docker Compose

```bash
# Production deployment
docker-compose up -d analytic-service
```

## Monitoring

The service includes a health check endpoint at `/api/v1/health` that can be used for:
- Kubernetes liveness/readiness probes
- Load balancer health checks
- Service monitoring and alerting

## Troubleshooting

### Port Already in Use
```bash
# Change PORT in .env file
PORT=8001

# Or kill the process using the port
lsof -ti:8000 | xargs kill -9
```

### Docker Issues
```bash
# Clean Docker resources
docker-compose down -v
docker system prune -a

# Rebuild from scratch
make docker-build
```

### Virtual Environment Issues
```bash
# Remove and recreate
rm -rf venv
make install
```

## Contributing

1. Follow the existing code style
2. Add tests for new features
3. Update documentation
4. Run all quality checks before committing

## License

This project is part of the Re:Path platform.

## Support

For issues and questions, please open an issue in the project repository.
