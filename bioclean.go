package timeseddel

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Bioclean is a timeseddel manager for Bioclean specific workshift tables
type Bioclean struct {
	rootpath string
}

const (
	lastDayInWorkshiftTable = 14
)

// NewBiocleanManager is a Bioclean struct constructor
func NewBiocleanManager(rootpath string) (Manager, error) {
	fileInfo, err := os.Stat(rootpath)

	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, errors.New("rootpath must point to a directory")
	}

	var manager Bioclean
	manager = Bioclean{
		rootpath: rootpath,
	}

	return &manager, nil
}

// GetWorkInfo tries to find the workinfo saved under the specified date
func (b Bioclean) GetWorkInfo(date time.Time) (WorkInfo, error) {
	panic("Not implemented")
}

// WriteWorkInfo writes the specified work info under the specified date
func (b Bioclean) WriteWorkInfo(date time.Time, workInfo WorkInfo) error {
	panic("Not implemented")
}

// getFilename gives the supposed filename based on the date passed.
// NOTE: This function doesn't check if the file exists,
// it just returns the string that's supposed to be the filename
func (b Bioclean) getFilename(date time.Time) string {
	filenameFormat := "%s Salary - %s %d.xlsx"
	year, month, day := date.Date()

	// increment the month if needed
	if day > lastDayInWorkshiftTable {
		if month == time.December {
			month = time.January
			year++
		} else {
			month++
		}
	}

	// add a leading 0 if needed
	var monthNum string
	if month < 10 {
		monthNum = "0" + strconv.Itoa(int(month))
	} else {
		monthNum = strconv.Itoa(int(month))
	}

	return fmt.Sprintf(filenameFormat, monthNum, month, year)
}
