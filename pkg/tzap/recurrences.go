// package tzap provides the implementation for Tzap which is responsible for
// handling recurrences and global occurrences.

package tzap

import (
	"strings"
	"sync"
)

// CheckAndHandleRecurrences counts the number of recurrences of the given filename
// within the data and calls either noReccurance or handleRecurrence based on the
// threshold provided as references.
func (t *Tzap) CheckAndHandleRecurrences(references int, filename string, noReccurance, handleRecurrence func(*Tzap) *Tzap) *Tzap {
	filepathCount := countMatchingFilepathValues(t, filename)
	if filepathCount < references {
		return noReccurance(t)
	}
	return handleRecurrence(t)
}

// countMatchingFilepathValues iterates through Tzap instances to count the number
// of occurrences of the given filepath value.
func countMatchingFilepathValues(t *Tzap, filepathValue interface{}) int {
	counter := 0
	for t != nil {
		if t.Data != nil {
			if value, ok := t.Data["filepath"]; ok && value == filepathValue {
				counter++
			}
		}
		t = t.Parent
	}
	return counter
}

// Package-level variables for global filepath occurrences tracking.
var (
	filepathOccurrences *sync.Map = &sync.Map{}
)

// ResetFilepathOccurrences clears the global filepath occurrences tracking.
func ResetFilepathOccurrences() {
	filepathOccurrences = &sync.Map{}
}

// countGlobalMatchingFilepathValues counts the global occurrences of the
// given filepath value.
func countGlobalMatchingFilepathValues(filepathValue string) int {
	counter := 0
	if value, ok := filepathOccurrences.Load(filepathValue); ok {
		counter = value.(int)
	}
	filepathOccurrences.Store(filepathValue, counter+1)
	return counter
}

// CheckAndHandleGlobalOccurrences checks and handles the global occurrences
// of the given filename within the data. Calls either noOccurrence or
// handleOccurrence based on the provided references.
func (t *Tzap) CheckAndHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(*Tzap) *Tzap) *Tzap {
	filepathCount := countGlobalMatchingFilepathValues(filename)
	parentLength := TotalLength(t)
	if parentLength < 2 || strings.Contains(filename, "model") || strings.Contains(filename, "search") {
		if filepathCount < references {
			return noOccurrence(t)
		}
	}
	return handleOccurrence(t)
}

// FileMustContainHandleGlobalOccurrences checks and handles the global occurrences
// of the provided filename within the data. Calls either noOccurrence or
// handleOccurrence based on the provided references.
func (t *Tzap) FileMustContainHandleGlobalOccurrences(references int, filename string, noOccurrence, handleOccurrence func(*Tzap) *Tzap) *Tzap {
	filepathCount := countGlobalMatchingFilepathValues(filename)
	if filepathCount < references {
		return noOccurrence(t)
	}
	return handleOccurrence(t)
}

// TotalLength calculates the total length of the Tzap instance chain.
func TotalLength(t *Tzap) int {
	length := 0
	for current := t; current != nil; current = current.Parent {
		length++
	}
	return length
}
