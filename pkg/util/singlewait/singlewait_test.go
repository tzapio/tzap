package singlewait_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/embed/localdb/singlewait"
)

func TestSingleWait(t *testing.T) {
	fn := func() string {
		return "Hello, World!"
	}

	s := singlewait.New(fn)

	data := s.GetData()

	if data != "Hello, World!" {
		t.Errorf("Unexpected data: got '%v' want '%v'", data, "Hello, World!")
	}

}
