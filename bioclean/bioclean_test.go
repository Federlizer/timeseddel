package bioclean

import (
	"testing"
	"time"

	"github.com/Federlizer/timeseddel"
)

const testdataPath = "./testdata/"

// Tests probably suck...

func TestNewBiocleanManager(t *testing.T) {
	t.Log("Constructs a new timeseddel manager...")
	manager, err := NewBiocleanManager(testdataPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Casts timeseddel manager to Bioclean manager...")
	var biocleanManager, ok = manager.(*Bioclean)
	if !ok {
		t.Fatal("Manager couldn't be cast to Bioclean")
	}

	if biocleanManager.rootpath != testdataPath {
		t.Log("Rootpath is different than passed path")
		t.Logf("Expected: %s\nActual: %s", testdataPath, biocleanManager.rootpath)
		t.Fail()
	}
}

func TestNewBiocleanManagerFailsOnRootpathFile(t *testing.T) {
	rootpath := testdataPath + "10 Salary - October 2019.xlsx"

	t.Log("Constructs a new timeseddel manager with wrong path...")
	_, err := NewBiocleanManager(rootpath)

	if err == nil {
		t.Fatal("Error shouldn't be nil")
	}
}

func TestNewBiocleanManagerFailsOnNonExistingRootpath(t *testing.T) {
	rootpath := "./somepath/that/does/not/exist"

	t.Log("Constructs a new timeseddel manager with wrong path...")
	_, err := NewBiocleanManager(rootpath)

	if err == nil {
		t.Fatal("Error shouldn't be nil")
	}
}

func TestGetWorkInfo(t *testing.T) {
	t.Log("Tests the normal behaviour of GetWorkInfo")

	dates := []time.Time{
		time.Date(2018, time.December, 17, 0, 0, 0, 0, time.UTC),
		time.Date(2019, time.January, 15, 0, 0, 0, 0, time.UTC),
	}

	expected := []*timeseddel.WorkInfo{
		&timeseddel.WorkInfo{
			Start: time.Date(2018, time.December, 17, 17, 0, 0, 0, time.UTC),
			End:   time.Date(2018, time.December, 17, 20, 0, 0, 0, time.UTC),
			Info:  "Picking up trash from 14, Cleaning reception in 14",
		},
		&timeseddel.WorkInfo{
			Start: time.Date(2019, time.January, 15, 17, 0, 0, 0, time.UTC),
			End:   time.Date(2019, time.January, 15, 20, 0, 0, 0, time.UTC),
			Info:  "Picking up trash in 14, cleaning main stairway in 46",
		},
	}

	t.Log("Instantiates a bioclean manager...")
	manager, err := NewBiocleanManager(testdataPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Checking method")
	for i := 0; i < len(dates); i++ {
		got, err := manager.GetWorkInfo(dates[i])
		if err != nil {
			t.Fatal(err)
		}

		if got.Start != expected[i].Start {
			t.Log("WorkInfo.Start")
			t.Logf("Expected: %v", expected[i].Start)
			t.Logf("Got     : %v", got.Start)
			t.Fail()
		}

		if got.End != expected[i].End {
			t.Log("WorkInfo.End")
			t.Logf("Expected: %v", expected[i].End)
			t.Logf("Got     : %v", got.End)
			t.Fail()
		}

		if got.Info != expected[i].Info {
			t.Log("WorkInfo.Info")
			t.Logf("Expected: %v", expected[i].Info)
			t.Logf("Got     : %v", got.Info)
			t.Fail()
		}
	}
}

func TestGetWorkInfoFailsOnMissingDate(t *testing.T) {
	t.Log("Tests if the GetWorkInfo function fails" +
		" when the date passed isn't held in any files")

	dates := []time.Time{
		time.Date(2017, time.December, 17, 0, 0, 0, 0, time.UTC),
		time.Date(2090, time.January, 15, 0, 0, 0, 0, time.UTC),
	}

	t.Log("Instantiates a bioclean manager...")
	manager, err := NewBiocleanManager(testdataPath)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(dates); i++ {
		got, err := manager.GetWorkInfo(dates[i])
		if err == nil {
			t.Log("When checking error\n")
			t.Log("Expected: error to be an Error object")
			t.Logf("Got (workinfo): %v", got)
			t.Logf("Got (err): %v", err)
			t.FailNow()
		}
	}
}

func TestWriteWorkingInfo(t *testing.T) {
	t.Log("Start writing info test")

	dates := []time.Time{
		time.Date(2019, time.October, 4, 0, 0, 0, 0, time.UTC),
	}

	workInfos := []timeseddel.WorkInfo{
		timeseddel.WorkInfo{
			Start: time.Date(2019, time.October, 4, 17, 0, 0, 0, time.UTC),
			End:   time.Date(2019, time.October, 4, 20, 0, 0, 0, time.UTC),
			Info:  "Pick up trash from 14, clean spots in the stairs in 14",
		},
	}

	t.Log("Creating Bioclean manager")
	manager, err := NewBiocleanManager(testdataPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Starting checks..")
	for i := 0; i < len(dates); i++ {
		date := dates[i]
		workInfo := workInfos[i]

		err := manager.WriteWorkInfo(date, workInfo)
		if err != nil {
			t.Log("Received an error when expecting no error")
			t.Fatal(err)
		}
	}
}
