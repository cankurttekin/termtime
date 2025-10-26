package model

import "time"

type CommandRecord struct {
	Command   string
	Timestamp time.Time
}
