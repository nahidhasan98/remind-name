package migration

import (
	"fmt"
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
		fmt.Printf("Running migration: %s\n", cmd.Name())
		if err := cmd.Execute(); err != nil {
			return fmt.Errorf("migration %s failed: %v", cmd.Name(), err)
		}
		fmt.Printf("Completed migration: %s\n", cmd.Name())
	}
	return nil
}
