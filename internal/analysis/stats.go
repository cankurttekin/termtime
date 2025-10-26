package analysis

import (
	"sort"
	"strings"
	"time"

	"termtime/internal/model"
)

type Statistics struct {
	CommandCounts map[string]int
	HourCounts    map[int]int
	DayCounts     map[string]int
	TimeSpan      TimeSpan
	HasTimestamps bool
}

type CommandStats struct {
	Command string
	Count   int
}

type HourStats struct {
	Hour  int
	Count int
}

type TimeSpan struct {
	First time.Time
	Last  time.Time
}

func Analyze(records []model.CommandRecord) *Statistics {
	stats := &Statistics{
		CommandCounts: make(map[string]int),
		HourCounts:    make(map[int]int),
		DayCounts:      make(map[string]int),
	}

	for _, r := range records {
		stats.processCommand(r)
		stats.processTimestamp(r)
	}

	return stats
}

func (s *Statistics) processCommand(record model.CommandRecord) {
	parts := strings.Fields(record.Command)
	if len(parts) > 0 {
		s.CommandCounts[parts[0]]++
	}
}

func (s *Statistics) processTimestamp(record model.CommandRecord) {
	if record.Timestamp.IsZero() {
		return
	}

	s.HasTimestamps = true
	s.HourCounts[record.Timestamp.Hour()]++
	s.DayCounts[record.Timestamp.Weekday().String()]++

	if s.TimeSpan.First.IsZero() || record.Timestamp.Before(s.TimeSpan.First) {
		s.TimeSpan.First = record.Timestamp
	}
	if record.Timestamp.After(s.TimeSpan.Last) {
		s.TimeSpan.Last = record.Timestamp
	}
}

func (s *Statistics) GetTopCommands(limit int) []CommandStats {
	var stats []CommandStats
	for cmd, count := range s.CommandCounts {
		stats = append(stats, CommandStats{Command: cmd, Count: count})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	if limit > 0 && limit < len(stats) {
		return stats[:limit]
	}
	return stats
}
