package cmdutil_test

import (
	"testing"

	"github.com/tzapio/tzap/cli/cmd/cmdutil"
)

func TestFilterWithIncludePattern(t *testing.T) {
	testCases := []struct {
		name            string
		inFiles         []string
		includePattern  []string
		expectedOutputs []string
	}{
		{
			name:            "Matching patterns",
			inFiles:         []string{"file1.txt", "file2.jpg", "file3.png", "file4.txt"},
			includePattern:  []string{"*.jpg", "*.png"},
			expectedOutputs: []string{"file2.jpg", "file3.png"},
		},
		{
			name:            "No matches",
			inFiles:         []string{"file1.txt", "file2.txt", "file3.txt"},
			includePattern:  []string{"*.jpg", "*.png"},
			expectedOutputs: []string{},
		},
		{
			name:            "Empty input",
			inFiles:         []string{},
			includePattern:  []string{"*.jpg", "*.png"},
			expectedOutputs: []string{},
		},
		{
			name:            "Empty include pattern",
			inFiles:         []string{"file1.txt", "file2.txt", "file3.txt"},
			includePattern:  []string{},
			expectedOutputs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outFiles := cmdutil.FilterWithIncludePattern(tc.inFiles, tc.includePattern)

			if len(outFiles) != len(tc.expectedOutputs) {
				t.Errorf("FilterWithIncludePattern returned incorrect output. Expected: %v, Got: %v", tc.expectedOutputs, outFiles)
				return
			}

			for i, file := range outFiles {
				expectedOutput := tc.expectedOutputs[i]
				if file != expectedOutput {
					t.Errorf("FilterWithIncludePattern returned incorrect output. Expected: %v, Got: %v", tc.expectedOutputs, outFiles)
					break
				}
			}
		})
	}
}
