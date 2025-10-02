.PHONY: help up down restart logs build clean dev-setup

# Default target
help:
	@echo "Available commands:"
	@echo "  make up          - Start all services"
	@echo "  make down        - Stop all services"
	@echo "  make restart     - Restart all services"
	@echo "  make logs        - View logs from all services"
	@echo "  make logs-url    - View logs from url-service"
	@echo "  make build       - Rebuild all services"
	@echo "  make clean       - Remove all containers, volumes, and images"
	@echo "  make dev-setup   - Initial setup for development"
	@echo "  make shell-url   - Enter url-service container shell"

# Start all services
up:
	docker compose up -d
	@echo "✅ All services started!"
	@echo "🌐 API Gateway: http://localhost:8080"
	@echo "🔗 URL Service: http://localhost:3001"

# Stop all services
down:
	docker compose down
	@echo "✅ All services stopped!"

# Restart all services
restart:
	docker compose restart
	@echo "✅ All services restarted!"

# View logs
logs:
	docker compose logs -f

logs-url:
	docker compose logs -f url-service

logs-nginx:
	docker compose logs -f nginx

# Rebuild services
build:
	docker compose build --no-cache
	@echo "✅ All services rebuilt!"

# Clean everything
clean:
	docker compose down -v --rmi all --remove-orphans
	@echo "✅ All containers, volumes, and images removed!"

# Development setup
dev-setup:
	@echo "🚀 Setting up development environment..."
	@mkdir -p nginx
	@if [ ! -f nginx/nginx.conf ]; then \
		echo "⚠️  nginx.conf not found. Creating default..."; \
	fi
	docker compose build
	docker compose up -d
	@echo "✅ Development environment ready!"

# Enter container shell
shell-url:
	docker compose exec url-service sh

shell-nginx:
	docker compose exec nginx sh

# Install dependencies in service
install-url:
	docker compose exec url-service npm install

# Run migrations (example)
migrate-url:
	docker compose exec url-service npm run migration:run

# Generate migration (example)
migration-generate-url:
	docker compose exec url-service npm run migration:generate -- -n $(name)

# Check service status
status:
	docker compose ps

# View resource usage
stats:
	docker stats