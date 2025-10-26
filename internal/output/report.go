package output

import (
	"fmt"

	"termtime/internal/analysis"
)

func PrintTopCommands(stats []analysis.CommandStats) {
	fmt.Println("=== Most Used Commands ===")
	for i, stat := range stats {
		fmt.Printf("%2d. %-10s %d\n", i+1, stat.Command, stat.Count)
	}
}

func PrintHourlyChart(hourCounts map[int]int) {
	fmt.Println("\n=== Hourly Activity ===")

	items := make([]ChartItem, 24)
	for i := range 24 {
		items[i] = ChartItem{
			Label: fmt.Sprintf("%02d:00", i),
			Value: hourCounts[i],
		}
	}

	DrawBarChart(items, 50)
}

func PrintDayOfWeekChart(dayCounts map[string]int) {
	fmt.Println("\n=== Days of Week ===")

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	items := make([]ChartItem, 0, 7)
	
	for _, day := range days {
		items = append(items, ChartItem{
			Label: day[:3],
			Value: dayCounts[day],
		})
	}

	DrawBarChart(items, 40)
}

func PrintTimeSpan(span analysis.TimeSpan) {
	fmt.Printf("From %s to %s\n", span.First.Format("02 Jan 2006 15:04"), span.Last.Format("02 Jan 2006 15:04"))
}

func PrintNoTimestamps() {
	fmt.Println("\nNo timestamps found — history doesn't include time info.")
}

