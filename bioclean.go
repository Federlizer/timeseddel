package timeseddel

import "time"

// Bioclean is a timeseddel manager for Bioclean specific workshift tables
type Bioclean struct {
	rootpath string
}

// GetWorkInfo tries to find the workinfo saved under the specified date
func (b Bioclean) GetWorkInfo(date time.Time) (WorkInfo, error) {
	panic("Not implemented")
}

// WriteWorkInfo writes the specified work info under the specified date
func (b Bioclean) WriteWorkInfo(date time.Time, workInfo WorkInfo) error {
	panic("Not implemented")
}
