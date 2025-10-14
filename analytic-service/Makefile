.PHONY: help install dev run test lint format type-check clean build docker-build docker-run docker-stop docker-logs

PYTHON := python3
VENV := venv
BIN := $(VENV)/bin
APP_NAME := analytic-service

help:
	@echo "🚀 Re:Path Analytics Service - Makefile Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make help         - Show this help message"
	@echo "  make install      - Create virtual environment and install dependencies"
	@echo "  make install-dev  - Install development dependencies"
	@echo "  make dev          - Run development server with hot reload"
	@echo "  make run          - Run production server"
	@echo "  make test         - Run tests with coverage"
	@echo "  make lint         - Run linting checks"
	@echo "  make format       - Format code with ruff"
	@echo "  make type-check   - Run type checking with mypy"
	@echo "  make clean        - Clean temporary files and caches"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make docker-dev   - Run Docker container with hot-reload"
	@echo "  make docker-stop  - Stop Docker containers"
	@echo "  make docker-logs  - View Docker container logs"
	@echo ""

install:
	@echo "📦 Creating virtual environment..."
	$(PYTHON) -m venv $(VENV)
	@echo "📦 Installing dependencies..."
	$(BIN)/pip install --upgrade pip
	$(BIN)/pip install -r requirements.txt
	@echo "✅ Installation complete!"
	@echo "💡 Activate virtual environment: source $(VENV)/bin/activate"

install-dev: install
	@echo "📦 Installing development dependencies..."
	$(BIN)/pip install -r requirements-dev.txt
	@echo "✅ Development dependencies installed!"

dev:
	@echo "🔥 Starting $(APP_NAME) in development mode with hot reload..."
	@if [ ! -d "$(VENV)" ]; then \
		echo "❌ Virtual environment not found. Run: make install"; \
		exit 1; \
	fi
	@if [ ! -f ".env" ]; then \
		echo "⚠️  .env file not found. Copying from .env.example..."; \
		cp .env.example .env; \
	fi
	$(BIN)/uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

run:
	@echo "🚀 Running $(APP_NAME) in production mode..."
	@if [ ! -d "$(VENV)" ]; then \
		echo "❌ Virtual environment not found. Run: make install"; \
		exit 1; \
	fi
	@if [ ! -f ".env" ]; then \
		echo "⚠️  .env file not found. Copying from .env.example..."; \
		cp .env.example .env; \
	fi
	$(BIN)/uvicorn app.main:app --host 0.0.0.0 --port 8000

test:
	@echo "🧪 Running tests with coverage..."
	@if [ ! -d "$(VENV)" ]; then \
		echo "❌ Virtual environment not found. Run: make install-dev"; \
		exit 1; \
	fi
	$(BIN)/pytest
	@echo "✅ Tests complete!"

lint:
	@echo "🔍 Running linting checks..."
	@if [ ! -d "$(VENV)" ]; then \
		echo "❌ Virtual environment not found. Run: make install-dev"; \
		exit 1; \
	fi
	$(BIN)/ruff check app tests
	@echo "✅ Linting complete!"

format:
	@echo "✨ Formatting code..."
	@if [ ! -d "$(VENV)" ]; then \
		echo "❌ Virtual environment not found. Run: make install-dev"; \
		exit 1; \
	fi
	$(BIN)/ruff check --fix app tests
	$(BIN)/ruff format app tests
	@echo "✅ Formatting complete!"

type-check:
	@echo "🔎 Running type checks..."
	@if [ ! -d "$(VENV)" ]; then \
		echo "❌ Virtual environment not found. Run: make install-dev"; \
		exit 1; \
	fi
	$(BIN)/mypy app
	@echo "✅ Type checking complete!"

clean:
	@echo "🧹 Cleaning temporary files and caches..."
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name ".pytest_cache" -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name ".ruff_cache" -exec rm -rf {} + 2>/dev/null || true
	find . -type d -name ".mypy_cache" -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete
	find . -type f -name "*.pyo" -delete
	find . -type f -name ".coverage" -delete
	rm -rf htmlcov coverage.html
	@echo "✅ Clean complete!"

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):latest .
	@echo "✅ Docker image built!"

docker-run:
	@echo "🐳 Running Docker container..."
	docker-compose up -d analytic-service
	@echo "✅ Container started!"
	@echo "📍 Access at: http://localhost:8000"
	@echo "📚 API Docs: http://localhost:8000/docs"

docker-dev:
	@echo "🐳 Running Docker container with hot-reload..."
	docker-compose --profile dev up -d analytic-service-dev
	@echo "✅ Development container started!"
	@echo "📍 Access at: http://localhost:8001"
	@echo "📚 API Docs: http://localhost:8001/docs"

docker-stop:
	@echo "🐳 Stopping Docker containers..."
	docker-compose down
	@echo "✅ Containers stopped!"

docker-logs:
	@echo "📜 Viewing Docker logs..."
	docker-compose logs -f analytic-service
