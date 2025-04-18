package main

import (
	"log"

	"github.com/nahidhasan98/remind-name/migration"
	"github.com/nahidhasan98/remind-name/migration/commands"
)

func main() {
	// Create a new migration runner
	runner := migration.NewRunner()

	// Add migration commands
	runner.AddCommand(commands.NewPlatformMigration())
	runner.AddCommand(commands.NewNameMigration())

	// Run all migrations
	if err := runner.Run(); err != nil {
		log.Fatal("Migration failed:", err)
	}
}
