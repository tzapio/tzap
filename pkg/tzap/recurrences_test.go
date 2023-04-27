package tzap_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

func Test_CheckAndHandleRecurrences_noRecurrences_expectNoRecurrenceExecution(t *testing.T) {
	root := tzap.InternalNew().
		AddTzap(&tzap.Tzap{Name: "tzap1", Data: types.MappedInterface{"filepath": "example.txt"}}).
		AddTzap(&tzap.Tzap{Name: "tzap2", Data: types.MappedInterface{"filepath": "example.txt"}}).
		AddTzap(&tzap.Tzap{Name: "tzap3", Data: types.MappedInterface{"filepath": "example.txt"}}).
		AddTzap(&tzap.Tzap{Name: "tzap4", Data: types.MappedInterface{"filepath": "example.txt"}}).
		AddTzap(&tzap.Tzap{Name: "tzap5"})

	noRecurrenceExecuted := false
	handleRecurrenceExecuted := false

	root.CheckAndHandleRecurrences(1, "example.txt",
		func(t *tzap.Tzap) *tzap.Tzap {
			noRecurrenceExecuted = true
			return t
		},
		func(t *tzap.Tzap) *tzap.Tzap {
			handleRecurrenceExecuted = true
			return t
		},
	)

	if noRecurrenceExecuted {
		t.Error("No recurrence function shouldnt be executed")
	}
	if !handleRecurrenceExecuted {
		t.Error("Handle recurrence function should be executed")
	}
}

func Test_CheckAndHandleRecurrences_oneRecurrence_expectHandleRecurrenceExecution(t *testing.T) {
	root := tzap.InternalNew()
	tzap1 := root.AddTzap(&tzap.Tzap{Name: "tzap1"})
	tzap2 := tzap1.AddTzap(&tzap.Tzap{Name: "tzap2", Data: types.MappedInterface{"filepath": "example.txt"}})
	tzap3 := tzap2.AddTzap(&tzap.Tzap{Name: "tzap3", Data: types.MappedInterface{"filepath": "example.txt"}})

	noRecurrenceExecuted := false
	handleRecurrenceExecuted := false

	tzap3.CheckAndHandleRecurrences(1, "example.txt",
		func(t *tzap.Tzap) *tzap.Tzap {
			noRecurrenceExecuted = true
			return t
		},
		func(t *tzap.Tzap) *tzap.Tzap {
			handleRecurrenceExecuted = true
			return t
		},
	)

	if noRecurrenceExecuted {
		t.Error("No recurrence function should not be executed")
	}
	if !handleRecurrenceExecuted {
		t.Error("Handle recurrence function should be executed")
	}
}

func Test_CheckAndHandleGlobalOccurrences_noOccurrences_expectNoOccurrenceExecution(t *testing.T) {
	root := tzap.InternalNew()
	defer tzap.ResetFilepathOccurrences()
	tzap1 := root.AddTzap(&tzap.Tzap{Name: "tzap1"})
	tzap2 := tzap1.AddTzap(&tzap.Tzap{Name: "tzap2", Data: types.MappedInterface{"filepath": "example.txt"}})
	tzap3 := tzap2.AddTzap(&tzap.Tzap{Name: "tzap3"})

	noOccurrenceExecuted := false
	handleOccurrenceExecuted := false

	tzap3.CheckAndHandleGlobalOccurrences(1, "example.txt",
		func(t *tzap.Tzap) *tzap.Tzap {
			noOccurrenceExecuted = true
			return t
		},
		func(t *tzap.Tzap) *tzap.Tzap {
			handleOccurrenceExecuted = true
			return t
		},
	)

	if noOccurrenceExecuted {
		t.Error("No occurrence function should be executed")
	}
	if !handleOccurrenceExecuted {
		t.Error("Handle occurrence function should not be executed")
	}
}

func Test_CheckAndHandleGlobalOccurrences_oneOccurrence_expectHandleOccurrenceExecution(t *testing.T) {
	defer tzap.ResetFilepathOccurrences()

	root := tzap.InternalNew()
	tzap1 := root.AddTzap(&tzap.Tzap{Name: "tzap1"})
	tzap2 := tzap1.AddTzap(&tzap.Tzap{Name: "tzap2", Data: types.MappedInterface{"filepath": "example.txt"}})
	tzap3 := tzap2.AddTzap(&tzap.Tzap{Name: "tzap3", Data: types.MappedInterface{"filepath": "example.txt"}})

	noOccurrenceExecuted := false
	handleOccurrenceExecuted := false

	tzap3.CheckAndHandleGlobalOccurrences(1, "example.txt",
		func(t *tzap.Tzap) *tzap.Tzap {
			noOccurrenceExecuted = true
			return t
		},
		func(t *tzap.Tzap) *tzap.Tzap {
			handleOccurrenceExecuted = true
			return t
		},
	)

	if noOccurrenceExecuted {
		t.Error("No occurrence function should not be executed")
	}
	if !handleOccurrenceExecuted {
		t.Error("Handle occurrence function should be executed")
	}
}

func Test_FileMustContainHandleGlobalOccurrences_noOccurrences_expectNoOccurrenceExecution(t *testing.T) {
	defer tzap.ResetFilepathOccurrences()

	root := tzap.InternalNew()
	tzap1 := root.AddTzap(&tzap.Tzap{Name: "tzap1"})
	tzap2 := tzap1.AddTzap(&tzap.Tzap{Name: "tzap2", Data: types.MappedInterface{"filepath": "example.txt"}})
	tzap3 := tzap2.AddTzap(&tzap.Tzap{Name: "tzap3"})

	noOccurrenceExecuted := false
	handleOccurrenceExecuted := false

	tzap3.FileMustContainHandleGlobalOccurrences(1, "example.txt",
		func(t *tzap.Tzap) *tzap.Tzap {
			noOccurrenceExecuted = true
			return t
		},
		func(t *tzap.Tzap) *tzap.Tzap {
			handleOccurrenceExecuted = true
			return t
		},
	)

	if !noOccurrenceExecuted {
		t.Error("No occurrence function should be executed")
	}
	if handleOccurrenceExecuted {
		t.Error("Handle occurrence function should not be executed")
	}
}

func Test_FileMustContainHandleGlobalOccurrences_oneOccurrence_expectHandleOccurrenceExecution(t *testing.T) {
	defer tzap.ResetFilepathOccurrences()

	root := tzap.InternalNew()
	tzap1 := root.AddTzap(&tzap.Tzap{Name: "tzap1"})
	tzap2 := tzap1.AddTzap(&tzap.Tzap{Name: "tzap2", Data: types.MappedInterface{"filepath": "example.txt"}})
	tzap3 := tzap2.AddTzap(&tzap.Tzap{Name: "tzap3", Data: types.MappedInterface{"filepath": "example.txt"}})

	noOccurrenceExecuted := false
	handleOccurrenceExecuted := false

	tzap3.FileMustContainHandleGlobalOccurrences(1, "example.txt",
		func(t *tzap.Tzap) *tzap.Tzap {
			noOccurrenceExecuted = true
			return t
		},
		func(t *tzap.Tzap) *tzap.Tzap {
			handleOccurrenceExecuted = true
			return t
		},
	)

	if !noOccurrenceExecuted {
		t.Error("No occurrence function should be executed")
	}
	if handleOccurrenceExecuted {
		t.Error("Handle occurrence function should not be executed")
	}
}
