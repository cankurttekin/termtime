package output

import (
	"fmt"
	"strings"
)

type ChartItem struct {
	Label string
	Value int
}

func DrawBarChart(items []ChartItem, barWidth int) {
	max := 0
	for _, item := range items {
		if item.Value > max {
			max = item.Value
		}
	}

	for _, item := range items {
		barLength := 0
		if max > 0 {
			barLength = int(float64(item.Value) / float64(max) * float64(barWidth))
		}

		bar := strings.Repeat("█", barLength)
		fmt.Printf("%s │%s %d\n", item.Label, bar, item.Value)
	}
}
