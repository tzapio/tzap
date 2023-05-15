package cmdutil_test

import (
	"testing"

	"github.com/tzapio/tzap/cli/cmd/cmdutil"
)

func TestFilterExcludePattern(t *testing.T) {
	testCases := []struct {
		name            string
		inFiles         []string
		excludePattern  []string
		expectedOutputs []string
	}{
		{
			name:            "Matching patterns",
			inFiles:         []string{"file1.txt", "file2.jpg", "file3.png", "file4.txt"},
			excludePattern:  []string{"*.jpg", "*.png"},
			expectedOutputs: []string{"file1.txt", "file4.txt"},
		},
		{
			name:            "No matches",
			inFiles:         []string{"file1.txt", "file2.txt", "file3.txt"},
			excludePattern:  []string{"*.jpg", "*.png"},
			expectedOutputs: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
		{
			name:            "Empty input",
			inFiles:         []string{},
			excludePattern:  []string{"*.jpg", "*.png"},
			expectedOutputs: []string{},
		},
		{
			name:            "Empty exclude pattern",
			inFiles:         []string{"file1.txt", "file2.txt", "file3.txt"},
			excludePattern:  []string{},
			expectedOutputs: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
		{
			name:            "Matching dir",
			inFiles:         []string{"file1.txt", "bar/world", "bar/world.txt", "hello/file1.txt", "world/hello/file2.txt", "world/world/hello/hello/file3.txt"},
			excludePattern:  []string{"hello"},
			expectedOutputs: []string{"file1.txt", "bar/world", "bar/world.txt"},
		},
		{
			name:            "Not Matching partial dir",
			inFiles:         []string{"file1.txt", "bar/world", "bar/world.txt", "hello2/file1.txt", "world/hello/file2.txt", "world/world/hello/hello/file3.txt"},
			excludePattern:  []string{"hello"},
			expectedOutputs: []string{"file1.txt", "bar/world", "bar/world.txt", "hello2/file1.txt"},
		},
		{
			name:            "Not Matching partial dir",
			inFiles:         []string{"file1.txt", "bar/world", "bar/world.txt", "hello2/file1.txt", "/hello/world", "world/hello/file2.txt", "world/world/hello/hello/file3.txt"},
			excludePattern:  []string{"hello"},
			expectedOutputs: []string{"file1.txt", "bar/world", "bar/world.txt", "hello2/file1.txt"},
		},
		{
			name:            "Matching dir/",
			inFiles:         []string{"file1.txt", "bar/world", "bar/world.txt", "hello2/file1.txt", "world/hello/file2.txt", "world/world/hello/hello/file3.txt"},
			excludePattern:  []string{"hello/"},
			expectedOutputs: []string{"file1.txt", "bar/world", "bar/world.txt", "hello2/file1.txt"},
		},
		{
			name:            "Matching absolute /dir ",
			inFiles:         []string{"file1.txt", "bar/world", "bar/world.txt", "hello/file1.txt", "world/hello/file2.txt", "world/world/hello/hello/file3.txt"},
			excludePattern:  []string{"/hello"},
			expectedOutputs: []string{"file1.txt", "bar/world", "bar/world.txt", "world/hello/file2.txt", "world/world/hello/hello/file3.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outFiles := cmdutil.FilterWithExcludePattern(tc.inFiles, tc.excludePattern)

			if len(outFiles) != len(tc.expectedOutputs) {
				t.Errorf("FilterExcludePattern returned incorrect output. Expected: %v, Got: %v", tc.expectedOutputs, outFiles)
				return
			}

			for i, file := range outFiles {
				expectedOutput := tc.expectedOutputs[i]
				if file != expectedOutput {
					t.Errorf("FilterExcludePattern returned incorrect output. Expected: %v, Got: %v", tc.expectedOutputs, outFiles)
					break
				}
			}
		})
	}
}
