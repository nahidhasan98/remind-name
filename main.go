package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	myApp "github.com/nahidhasan98/remind-name/app"
	"github.com/nahidhasan98/remind-name/bot"
	"github.com/nahidhasan98/remind-name/notification"
)

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
	log.Println("Shutting down gracefully...")

	// Cancel the context to signal goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("Application stopped.")
}

func main() {
	fmt.Println("program is running...")

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
