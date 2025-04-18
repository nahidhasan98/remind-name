package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nahidhasan98/remind-name/config"
)

type App struct {
	*gin.Engine
	server *http.Server // to support graceful shutdown
}

func New() *App {
	if config.APP_MODE == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &App{
		Engine: gin.Default(),
	}
}

func (as *App) Start(ctx context.Context) {
	// Create a new HTTP server with the Gin engine as the handler
	as.server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", config.APP_PORT),
		Handler: as.Engine,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("RUNNING: Web server on port %d.", config.APP_PORT)
		if err := as.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Web server failed: %v", err)
		}
	}()

	// Wait for the context to be canceled
	<-ctx.Done()

	log.Println("Shutting down web server...")

	// Create a shutdown context with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shut down the server
	if err := as.server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Web server shutdown failed: %v", err)
	}

	log.Println("Web server stopped.")
}
