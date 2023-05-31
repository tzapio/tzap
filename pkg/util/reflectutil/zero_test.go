package reflectutil

import (
	"reflect"
	"testing"
)

func TestIsZero(t *testing.T) {
	var tests = []struct {
		input interface{}
		want  bool
	}{
		{0, true},
		{1, false},
		{"", true},
		{"non-empty string", false},
		{nil, true},
		{[]int{}, false},
		{[]int{1, 2, 3}, false},
		{map[string]int{}, false},
		{map[string]int{"foo": 1}, false},
		{struct{}{}, true},
		{struct{ Name string }{"John"}, false},
	}

	for _, tt := range tests {
		reflectType := reflect.TypeOf(tt.input)
		var testname string
		if reflectType == nil {
			testname = "nil"
		} else {
			testname = reflectType.String()
		}
		t.Run(testname, func(t *testing.T) {
			ans := IsZero(tt.input)
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
