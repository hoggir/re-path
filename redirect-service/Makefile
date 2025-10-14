.PHONY: help install wire swagger build run dev test clean

BINARY_NAME=redirect-service
CMD_DIR=cmd
WIRE_DIR=cmd

help:
	@echo "🚀 Re:Path Redirect Service - Makefile Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make help      - Show this help message"
	@echo "  make install   - Install dependencies and dev tools"
	@echo "  make wire      - Generate wire dependency injection"
	@echo "  make swagger   - Generate Swagger documentation"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the built application"
	@echo "  make dev       - Run with hot reload (requires air)"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo ""

install:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "📦 Installing dev tools (air + wire + swag)..."
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Installation complete (Wire v0.6.0 - stable)"

## Generate dependency injection dengan Wire
wire:
	@echo "⚙️  Generating dependency injection with Wire..."
	cd $(WIRE_DIR) && wire
	@echo "✅ Wire generation complete!"

## Generate Swagger documentation
swagger:
	@echo "📚 Generating Swagger documentation..."
	swag init -g $(CMD_DIR)/main.go -o docs
	@echo "✅ Swagger docs generated at docs/"
	@echo "   Access at: http://localhost:3011/swagger/index.html"

build:
	@echo "🏗️  Building $(BINARY_NAME)..."
	@mkdir -p tmp
	go build -o tmp/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Build complete: tmp/$(BINARY_NAME)"

run:
	@echo "🚀 Running $(BINARY_NAME)..."
	./tmp/$(BINARY_NAME)

dev:
	@echo "🔥 Starting $(BINARY_NAME) in development mode with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ Air not installed. Run: make install"; \
		exit 1; \
	fi

test:
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests complete. Coverage report: coverage.html"

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf tmp/
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete"