# ===========================================
# 🧱 Go Redirect Service — Makefile
# ===========================================

# Nama binary yang dihasilkan
BINARY_NAME=redirect-service

# Lokasi direktori utama aplikasi
CMD_DIR=cmd/server

# Lokasi file wire.go
WIRE_DIR=cmd/server

# ===========================================
# 🧩 Perintah Development
# ===========================================

## Jalankan aplikasi dengan hot reload (pakai Air)
dev:
	@echo "🚀 Starting ${BINARY_NAME} in development mode..."
	@air

## Generate dependency injection dengan Wire
wire:
	@echo "⚙️  Generating dependency injection with Wire..."
	cd $(WIRE_DIR) && go run github.com/google/wire/cmd/wire
	@echo "✅ Wire generation complete!"

## Install Air (hot reload) dan Wire (dependency injection)
setup:
	@echo "📦 Installing dev tools (air + wire)..."
	go install github.com/cosmtrek/air@latest
	go get github.com/google/wire/cmd/wire@v0.7.0
	@echo "✅ Setup complete!"

# ===========================================
# 🏗️ Build & Run
# ===========================================

## Build aplikasi (tanpa hot reload)
build:
	@echo "🏗️  Building ${BINARY_NAME}..."
	go build -o tmp/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Build complete: tmp/$(BINARY_NAME)"

## Jalankan hasil build
run:
	@echo "🚀 Running ${BINARY_NAME}..."
	./tmp/$(BINARY_NAME)

## Bersihkan file hasil build
clean:
	@echo "🧹 Cleaning up build files..."
	rm -rf tmp
	@echo "✅ Cleanup done."
