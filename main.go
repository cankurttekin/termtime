package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"termtime/internal/analysis"
	"termtime/internal/config"
	"termtime/internal/history"
	"termtime/internal/model"
	"termtime/internal/output"
)

func main() {
	cfg := config.DefaultConfig()
	
	// Parse flags
	ignoreFlag := flag.String("ignore", "", "Comma-separated list of commands to ignore (e.g., 'cd,ls,pwd')")
	limitFlag := flag.Int("limit", 10, "Number of top commands to show")
	flag.Parse()
	
	// Apply flags to config
	cfg.TopCommandsLimit = *limitFlag
	cfg.SetIgnoreCommands(*ignoreFlag)
	
	// Validate config
	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid configuration:", err)
	}
	
	// Find history file
	histFile, shellType, err := history.FindHistoryFile()
	if err != nil {
		log.Fatal("Failed to find history file:", err)
	}

	// Parse history
	records, err := history.ParseFile(histFile)
	if err != nil {
		log.Fatal("Failed to parse history:", err)
	}

	if len(records) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: No command records found in history file")
		os.Exit(0)
	}
	
	// Filter ignored commands
	filteredRecords := filterRecords(records, cfg)
	
	if len(filteredRecords) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: All commands filtered by ignore list")
		os.Exit(0)
	}

	// Analyze records
	stats := analysis.Analyze(filteredRecords)
	
	// Print results
	fmt.Printf("Using %s history file: %s\n\n", shellType, histFile)
	if len(cfg.IgnoreCommands) > 0 {
		fmt.Printf("Ignoring commands: %v\n\n", cfg.IgnoreCommands)
	}
	
	// Print report
	output.PrintTopCommands(stats.GetTopCommands(cfg.TopCommandsLimit))
	
	if stats.HasTimestamps {
		output.PrintHourlyChart(stats.HourCounts)
		output.PrintDayOfWeekChart(stats.DayCounts)
		//fmt.Printf("\n")
		//output.PrintTimeSpan(stats.TimeSpan)
	} else {
		output.PrintNoTimestamps()
	}
}

func filterRecords(records []model.CommandRecord, cfg *config.Config) []model.CommandRecord {
	var filtered []model.CommandRecord
	for _, record := range records {
		if !cfg.ShouldIgnore(record.Command) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}
