package bioclean

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Federlizer/timeseddel"
)

// Bioclean is a timeseddel manager for Bioclean specific workshift tables
type Bioclean struct {
	rootpath string
}

// Bioclean config
const (
	lastDayInWorkshiftTable = 14

	sheetName = "Timeseddel"

	firstRow = 11

	endOfWeekSkipCount = 4

	dayOfTheWeekColumn = "B"
	startColumn        = "C"
	endColumn          = "D"
	infoColumn         = "F"
)

// NewBiocleanManager is a Bioclean struct constructor
func NewBiocleanManager(rootpath string) (timeseddel.Manager, error) {
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
func (b *Bioclean) GetWorkInfo(date time.Time) (*timeseddel.WorkInfo, error) {
	filename := b.getFilename(date)
	file, err := b.getFile(filename)

	if err != nil {
		return nil, err
	}

	row, err := b.getRow(file, date)
	if err != nil {
		panic(err)
	}

	startStr, err := file.GetCellValue(sheetName, startColumn+strconv.Itoa(row))
	if err != nil {
		panic(err)
	}

	endStr, err := file.GetCellValue(sheetName, endColumn+strconv.Itoa(row))
	if err != nil {
		panic(err)
	}

	info, err := file.GetCellValue(sheetName, infoColumn+strconv.Itoa(row))
	if err != nil {
		panic(err)
	}

	startTime, err := mapStringToTime(startStr)
	if err != nil {
		panic(err)
	}

	endTime, err := mapStringToTime(endStr)
	if err != nil {
		panic(err)
	}

	workinfo := &timeseddel.WorkInfo{
		Start: time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC),
		End:   time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC),
		Info:  info,
	}

	return workinfo, nil
}

// WriteWorkInfo writes the specified work info under the specified date
//
// CAUTION: since the function uses a lot of the utils funcitons which are tested in GetWorkInfo
// this function is not thoroughly tested, and the tests don't check wether or not
// the info has been written.
func (b *Bioclean) WriteWorkInfo(date time.Time, workInfo timeseddel.WorkInfo) error {
	filename := b.getFilename(date)
	file, err := b.getFile(filename)

	if err != nil {
		return err
	}

	row, err := b.getRow(file, date)
	if err != nil {
		return err
	}

	startValue := mapTimeToExcelFloat(workInfo.Start)
	fmt.Println(startValue)

	endValue := mapTimeToExcelFloat(workInfo.End)
	fmt.Println(endValue)

	startAxis := startColumn + strconv.Itoa(row)
	endAxis := endColumn + strconv.Itoa(row)
	infoAxis := infoColumn + strconv.Itoa(row)

	err = file.SetCellValue(sheetName, startAxis, startValue)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, endAxis, endValue)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, infoAxis, workInfo.Info)
	if err != nil {
		return err
	}

	err = file.Save()
	return err
}

// getFilename gives the supposed filename based on the date passed.
// NOTE: This function doesn't check if the file exists,
// it just returns the string that's supposed to be the filename
func (b *Bioclean) getFilename(date time.Time) string {
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

// getFile tries to find the file which corresponds to the passed filename
func (b *Bioclean) getFile(filename string) (*excelize.File, error) {
	var filePath string

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == filename {
			filePath = path
			return io.EOF
		}

		return nil
	}

	err := filepath.Walk(b.rootpath, walkFunc)
	if err != nil && err != io.EOF {
		return nil, err
	}

	file, err := excelize.OpenFile(filePath)
	return file, err
}

// getRow tries to find the row that corresponds to the passed date
func (b *Bioclean) getRow(file *excelize.File, date time.Time) (int, error) {
	currentRow := firstRow
	wantedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var currentRowDate time.Time
	day, _ := time.ParseDuration("24h")

	// Ajust to bioclean's table standards
	if wantedDate.Day() < 15 {
		if wantedDate.Month() == time.January {
			currentRowDate = time.Date(wantedDate.Year()-1, time.December, 15, 0, 0, 0, 0, time.UTC)
		} else {
			currentRowDate = time.Date(wantedDate.Year(), (wantedDate.Month() - 1), 15, 0, 0, 0, 0, time.UTC)
		}
	} else {
		currentRowDate = time.Date(wantedDate.Year(), wantedDate.Month(), 15, 0, 0, 0, 0, time.UTC)
	}

	for currentRowDate.Before(date) {
		cellAddress := dayOfTheWeekColumn + strconv.Itoa(currentRow)
		value, err := file.GetCellValue(sheetName, cellAddress)
		if err != nil {
			return -1, err
		}

		currentRowDate = currentRowDate.Add(day)
		if value == "Sunday" {
			currentRow += endOfWeekSkipCount
		} else {
			currentRow++
		}
	}

	return currentRow, nil
}

// HH:MM[:SS]
// mapStringToTime takes a string in the above format
// and returns a time that holds the hours minutes and, if found in the string, seconds
func mapStringToTime(str string) (time.Time, error) {
	data := strings.Split(str, ":")

	hours, err := strconv.Atoi(data[0])
	if err != nil {
		return time.Time{}, err
	}

	minutes, err := strconv.Atoi(data[1])
	if err != nil {
		return time.Time{}, err
	}

	seconds := 0
	if len(data) > 2 {
		seconds, err = strconv.Atoi(data[2])
		if err != nil {
			return time.Time{}, err
		}
	}

	date := time.Date(1, time.January, 1, hours, minutes, seconds, 0, time.UTC)
	return date, nil
}

// mapTimeToExcelFloat transforms the time passed into an
// excel consumable float64 for date and time
//
// please refer to (https://www.myonlinetraininghub.com/excel-date-and-time)
func mapTimeToExcelFloat(time time.Time) float64 {
	var value float64

	value += float64(time.Hour()) / 24     // add the hours
	value += float64(time.Minute()) / 1440 // add the minutes

	return value
}
