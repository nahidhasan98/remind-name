package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	myApp "github.com/nahidhasan98/remind-name/app"
	"github.com/nahidhasan98/remind-name/bot"
	"github.com/nahidhasan98/remind-name/config"
	"github.com/nahidhasan98/remind-name/logger"
	"github.com/nahidhasan98/remind-name/notification"
)

// Initialize the logger
func initializeLogger() {
	logger.Init(config.DEBUG_MODE, config.LOG_FILE)
	defer func() {
		if !config.DEBUG_MODE {
			logger.Info("Logger closed.")
		}
	}()
}

// Initialize all bots
func initializeBots(ctx context.Context, wg *sync.WaitGroup) {
	manager := bot.GetBotManager()

	wg.Add(1)
	go func() {
		defer wg.Done()
		manager.StartAll(ctx)
	}()
}

// Start the web server
func startWebServer(ctx context.Context, wg *sync.WaitGroup) *myApp.App {
	app := myApp.New()
	app.RegisterRoute()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Start(ctx)
	}()

	return app
}

// Start the notification scheduler with context & waitgroup
func startNotificationScheduler(ctx context.Context, wg *sync.WaitGroup) {
	notificationService := notification.NewNotificationService()

	wg.Add(1)
	go func() {
		defer wg.Done()
		notificationService.StartScheduler(ctx)
	}()
}

// Wait for termination signals and handle graceful shutdown
func waitForShutdown(cancel context.CancelFunc, wg *sync.WaitGroup) {
	// Wait for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Info("Shutting down gracefully...")

	// Cancel the context to signal goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info("Application stopped.")
}

func main() {
	// Initialize logger
	initializeLogger()
	logger.Info("Starting Remind Name application...")

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a wait group for graceful shutdown
	var wg sync.WaitGroup

	// Initialize and start bots
	initializeBots(ctx, &wg)

	// Start the web server
	startWebServer(ctx, &wg)

	// Start the notification scheduler
	startNotificationScheduler(ctx, &wg)

	// Handle shutdown signals
	waitForShutdown(cancel, &wg)
}
