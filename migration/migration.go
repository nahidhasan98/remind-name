package migration

import (
	"fmt"

	"github.com/nahidhasan98/remind-name/logger"
)

// Command interface for different migration commands
type Command interface {
	Execute() error
	Name() string
}

// Runner executes migration commands
type Runner struct {
	commands []Command
}

// NewRunner creates a new migration runner
func NewRunner() *Runner {
	return &Runner{
		commands: make([]Command, 0),
	}
}

// AddCommand adds a migration command to the runner
func (r *Runner) AddCommand(cmd Command) {
	r.commands = append(r.commands, cmd)
}

// Run executes all migration commands
func (r *Runner) Run() error {
	for _, cmd := range r.commands {
		logger.Info("Running migration: %s", cmd.Name())
		if err := cmd.Execute(); err != nil {
			logger.Error("Migration %s failed: %v", cmd.Name(), err)
			return fmt.Errorf("migration %s failed: %v", cmd.Name(), err)
		}
		logger.Info("Completed migration: %s", cmd.Name())
	}
	return nil
}
