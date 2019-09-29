package timeseddel

import "time"

// Manager is an instance that manages timeseddel workshift tables
type Manager interface {
	// GetWorkInfo tries to find the workinfo saved under the specified date
	GetWorkInfo(date time.Time) (WorkInfo, error)

	// WriteWorkInfo writes the specified work info under the specified date
	WriteWorkInfo(date time.Time, workinfo WorkInfo) error
}

// WorkInfo contains a work info for a specific date
type WorkInfo struct {
	start time.Time
	end   time.Time
	info  string
}
