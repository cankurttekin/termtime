package history

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ShellType represents the type of shell
type ShellType string

const (
	ShellZsh  ShellType = "zsh"
	ShellBash ShellType = "bash"
	ShellUnknown ShellType = "unknown"
)

// FindHistoryFile locates the shell history file and determines the shell type
func FindHistoryFile() (string, ShellType, error) {
	usr, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("failed to get current user: %w", err)
	}

	homeDir := usr.HomeDir
	
	// Detect current shell from environment
	currentShell := os.Getenv("SHELL")
	if currentShell == "" {
		currentShell = os.Getenv("SHELL")
	}
	
	shellType := detectShellType(currentShell)
	
	// Map shell types to their history files
	historyFiles := map[ShellType]string{
		ShellZsh:  filepath.Join(homeDir, ".zsh_history"),
		ShellBash: filepath.Join(homeDir, ".bash_history"),
	}
	
	// Try current shell first
	if hist := historyFiles[shellType]; hist != "" {
		if _, err := os.Stat(hist); err == nil {
			return hist, shellType, nil
		}
	}
	
	// Fallback: check all history files
	for st, hist := range historyFiles {
		if _, err := os.Stat(hist); err == nil {
			return hist, st, nil
		}
	}
	
	return "", ShellUnknown, fmt.Errorf("no bash or zsh history file found")
}

// detectShellType determines shell type from the shell path
func detectShellType(shellPath string) ShellType {
	shellPath = strings.ToLower(shellPath)
	
	if strings.Contains(shellPath, "zsh") {
		return ShellZsh
	}
	if strings.Contains(shellPath, "bash") {
		return ShellBash
	}
	
	return ShellUnknown
}

