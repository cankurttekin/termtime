package history

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"termtime/internal/model"
)

// Parser defines the interface for parsing shell history files
type Parser interface {
	Parse(file string) ([]model.CommandRecord, error)
}

// NewParser creates a parser for the given shell type
func NewParser(shellType string) (Parser, error) {
	switch shellType {
	case "zsh":
		return &zshParser{}, nil
	case "bash":
		return &bashParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported shell type: %s", shellType)
	}
}

// ParseFile uses the appropriate parser based on the file path
func ParseFile(file string) ([]model.CommandRecord, error) {
	parser := DetermineParser(file)
	return parser.Parse(file)
}

// DetermineParser returns the appropriate parser based on file path
func DetermineParser(file string) Parser {
	if strings.Contains(file, "zsh") {
		return &zshParser{}
	}
	return &bashParser{}
}

// zshParser handles zsh history format: ": 1698247289:0;ls -l"
type zshParser struct{}

func (p *zshParser) Parse(file string) ([]model.CommandRecord, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []model.CommandRecord
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if record := p.parseLine(line); record != nil {
			records = append(records, *record)
		}
	}

	return records, scanner.Err()
}

func (p *zshParser) parseLine(line string) *model.CommandRecord {
	line = strings.TrimSpace(line)
	if line == "" || !strings.HasPrefix(line, ": ") {
		return nil
	}

	// Format: ": 1698247289:0;ls -l"
	parts := strings.SplitN(line, ";", 2)
	if len(parts) != 2 {
		return nil
	}

	// Extract timestamp from meta part ": 1698247289:0"
	meta := strings.Fields(parts[0])
	if len(meta) < 1 {
		return nil
	}

	epochStr := strings.TrimPrefix(meta[0], ":")
	epochStr = strings.TrimSpace(epochStr)
	
	epoch, err := parseEpoch(epochStr)
	if err != nil {
		return nil
	}

	cmd := strings.TrimSpace(parts[1])
	return &model.CommandRecord{Command: cmd, Timestamp: epoch}
}

// bashParser handles bash history format
// Format with timestamps: "#1761504065\ncommand\n#1761504087\ncommand"
// Format without timestamps: just "command"
type bashParser struct{}

func (p *bashParser) Parse(file string) ([]model.CommandRecord, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []model.CommandRecord
	scanner := bufio.NewScanner(f)
	
	var timestamp time.Time
	hasTimestamp := false

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		
		if line == "" {
			continue
		}
		
		// Check if this is a timestamp line (starts with # followed by digits)
		if strings.HasPrefix(line, "#") {
			epochStr := strings.TrimPrefix(line, "#")
			if epoch, err := parseEpoch(epochStr); err == nil {
				timestamp = epoch
				hasTimestamp = true
				continue
			}
		}
		
		// This is a command line
		if hasTimestamp && !timestamp.IsZero() {
			records = append(records, model.CommandRecord{
				Command:   line,
				Timestamp: timestamp,
			})
			hasTimestamp = false // Reset for next command
		} else {
			// Bash history without timestamps
			records = append(records, model.CommandRecord{Command: line})
		}
	}

	return records, scanner.Err()
}

// parseEpoch converts a Unix timestamp string to time.Time
func parseEpoch(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}
