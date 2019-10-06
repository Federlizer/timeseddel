package timeseddel

import "time"

// WorkInfo contains a work info for a specific date
type WorkInfo struct {
	Start time.Time
	End   time.Time
	Info  string
}
