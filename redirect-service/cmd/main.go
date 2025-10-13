package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Re:Path Redirect Service API
// @version         1.0
// @description     A high-performance URL shortening and redirect service built with Go, Gin, MongoDB, and Redis.
// @description     This service provides URL shortening, custom aliases, and redirect functionality.

// @contact.name   API Support
// @contact.email  support@repath.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:3011
// @BasePath  /

// @schemes   http https

func main() {
	srv, err := InitializeApp()
	if err != nil {
		log.Fatalf("❌ Failed to initialize app: %v", err)
	}
	defer func() {
		log.Println("🧹 Cleaning up resources...")
		if err := srv.MongoDB.Close(); err != nil {
			log.Printf("Error closing MongoDB: %v", err)
		}
		if err := srv.Redis.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3011"
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: srv.GetRouter(),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("🚀 Redirect service starting on port %s...", port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		log.Fatalf("❌ Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("\n⚠️  Received signal: %v, starting graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("❌ Error during shutdown: %v", err)
			if err := httpServer.Close(); err != nil {
				log.Fatalf("❌ Could not stop server: %v", err)
			}
		}

		log.Println("✅ Server stopped successfully")
		log.Println("👋 Shutdown complete")
	}
}
