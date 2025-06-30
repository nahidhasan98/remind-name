package main

import (
	"github.com/nahidhasan98/remind-name/config"
	"github.com/nahidhasan98/remind-name/logger"

	"os"

	"github.com/nahidhasan98/remind-name/migration"
	"github.com/nahidhasan98/remind-name/migration/commands"
)

func main() {
	logger.Init(config.DEBUG_MODE, "logs/app.log")

	// Create a new migration runner
	runner := migration.NewRunner()

	// Add migration commands
	runner.AddCommand(commands.NewPlatformMigration())
	runner.AddCommand(commands.NewNameMigration())

	// Run all migrations
	if err := runner.Run(); err != nil {
		logger.Error("Migration failed: %v", err)
		os.Exit(1)
	}
}
