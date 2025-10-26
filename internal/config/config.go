package config

import (
	"fmt"
	"strings"
)

// Config holds configuration options for the analyzer
type Config struct {
	TopCommandsLimit int
	ShellType        string
	HistoryFile      string
	IgnoreCommands   []string
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		TopCommandsLimit: 10,
		ShellType:        "", // Auto-detect
		HistoryFile:      "", // Auto-detect
		IgnoreCommands:   []string{},
	}
}

// SetIgnoreCommands parses a comma-separated string of commands to ignore
func (c *Config) SetIgnoreCommands(ignoreStr string) {
	if ignoreStr == "" {
		return
	}
	
	commands := strings.Split(ignoreStr, ",")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" {
			c.IgnoreCommands = append(c.IgnoreCommands, cmd)
		}
	}
}

// ShouldIgnore checks if a command should be ignored
func (c *Config) ShouldIgnore(command string) bool {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}
	
	cmd := parts[0]
	for _, ignoreCmd := range c.IgnoreCommands {
		if cmd == ignoreCmd {
			return true
		}
	}
	return false
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.TopCommandsLimit < 0 {
		return fmt.Errorf("top commands limit cannot be negative")
	}
	return nil
}

