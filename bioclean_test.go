package timeseddel_test

import (
	"testing"

	"github.com/Federlizer/timeseddel"
)

func TestImplementsManager(t *testing.T) {
	var _ timeseddel.Manager = timeseddel.Bioclean{}
	var _ timeseddel.Manager = (*timeseddel.Bioclean)(nil)
}
