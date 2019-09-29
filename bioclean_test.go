package timeseddel

import (
	"testing"
	"time"
)

func TestImplementsManager(t *testing.T) {
	var _ Manager = Bioclean{}
	var _ Manager = (*Bioclean)(nil)
}

func TestConstructor(t *testing.T) {
	var rootpath string = "/home/federlizer/Dropbox/Nikola Velichkov"

	var manager Manager
	manager, err := NewBiocleanManager(rootpath)

	var biocleanManager, ok = manager.(*Bioclean)

	if !ok {
		t.Log("Manager couldn't be cast to Bioclean")
		t.Fail()
	}

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if biocleanManager.rootpath != rootpath {
		t.Log("Rootpath is different than passed path")
		t.Logf("Expected: %s\nActual: %s", rootpath, biocleanManager.rootpath)
		t.Fail()
	}
}

func TestConstructorFailsOnRootpathFile(t *testing.T) {
	rootpath := "./tests/testdata/10 Salary - October 2019.xlsx"

	var manager Manager
	manager, err := NewBiocleanManager(rootpath)

	if err == nil {
		t.Log("Error is nil")
		t.Logf("Error: %s", err)
		t.FailNow()
	}

	if manager != nil {
		t.Log("Manager isn't nil")
		t.Log(manager)
		t.FailNow()
	}
}

func TestConstructorFailsOnNonExistingRootpath(t *testing.T) {
	rootpath := "./somepath/that/does/not/exist"

	var manager Manager
	manager, err := NewBiocleanManager(rootpath)

	if err == nil {
		t.Log("Error is nil")
		t.Logf("Error: %s", err)
		t.FailNow()
	}

	if manager != nil {
		t.Log("Manager isn't nil")
		t.Log(manager)
		t.FailNow()
	}
}

func TestGetFilename(t *testing.T) {
	rootpath := "./tests/testdata"

	manager, err := NewBiocleanManager(rootpath)

	if err != nil {
		t.Fatal(err)
	}

	if manager == nil {
		t.Fatal("Manager is nil")
	}

	biocleanManager, ok := manager.(*Bioclean)
	if !ok {
		t.Fatal("Manager can't be cast to Bioclean")
	}

	date := time.Date(2019, time.February, 10, 0, 0, 0, 0, time.UTC)

	expectedFilename := "02 Salary - February 2019.xlsx"
	filename := biocleanManager.getFilename(date)
	if filename != expectedFilename {
		t.Logf("Expected: %v\nGot: %v", expectedFilename, filename)
		t.FailNow()
	}
}

func TestGetFilenameNextMonthTable(t *testing.T) {
	rootpath := "./tests/testdata"

	manager, err := NewBiocleanManager(rootpath)
	if err != nil {
		t.Fatal(err)
	}
	if manager == nil {
		t.Fatal("Manager is nil")
	}

	biocleanManager, ok := manager.(*Bioclean)
	if !ok {
		t.Fatal("Manager can't be cast to Bioclean")
	}

	date := time.Date(2019, time.February, 15, 0, 0, 0, 0, time.UTC)
	expectedFilename := "03 Salary - March 2019.xlsx"
	filename := biocleanManager.getFilename(date)
	if filename != expectedFilename {
		t.Logf("\nPassed: %v\nExpected: %v\nGot: %v", date, expectedFilename, filename)
		t.FailNow()
	}

	date = time.Date(2019, time.September, 29, 0, 0, 0, 0, time.UTC)
	expectedFilename = "10 Salary - October 2019.xlsx"
	filename = biocleanManager.getFilename(date)
	if filename != expectedFilename {
		t.Logf("\nPassed: %v\nExpected: %v\nGot: %v", date, expectedFilename, filename)
		t.FailNow()
	}

	date = time.Date(2018, time.December, 31, 0, 0, 0, 0, time.UTC)
	expectedFilename = "01 Salary - January 2019.xlsx"
	filename = biocleanManager.getFilename(date)
	if filename != expectedFilename {
		t.Logf("\nPassed: %v\nExpected: %v\nGot: %v", date, expectedFilename, filename)
		t.FailNow()
	}

	date = time.Date(2019, time.January, 15, 0, 0, 0, 0, time.UTC)
	expectedFilename = "02 Salary - February 2019.xlsx"
	filename = biocleanManager.getFilename(date)
	if filename != expectedFilename {
		t.Logf("\nPassed: %v\nExpected: %v\nGot: %v", date, expectedFilename, filename)
		t.FailNow()
	}
}
