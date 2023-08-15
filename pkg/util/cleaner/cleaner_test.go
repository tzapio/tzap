package cleaner_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/util/cleaner"
)

func TestFileWriteClean(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "Content without #### as first line",
			input:          "Content without #### as first line",
			expectedOutput: "Content without #### as first line",
		},
		{
			name:           "Content with #### as first line",
			input:          "####Content with \n#### as first line",
			expectedOutput: "#### as first line",
		},
		{
			name:           "Content with #### as first line and two code blocks, should return original content",
			input:          "####Content with \n#### as first line\n```go\nCode Block\n```, or this exists \n```go\nCode Block2\n``` what",
			expectedOutput: "#### as first line\n```go\nCode Block\n```, or this exists \n```go\nCode Block2\n``` what",
		},
		{
			name:           "Content with #### as first line and a code block should return the code block",
			input:          "####hello world\n```go\npackage integrations\nhello\n```\n, or this exists",
			expectedOutput: "package integrations\nhello\n",
		},
		{
			name:           "Content with #### as first line and a code block should return the code block",
			input:          "####hello world\n```go\npackage integrations\nhello\n```\n```\n, or this exists",
			expectedOutput: "package integrations\nhello\n",
		},
		{
			name:           "Content with #### as first line and a code block should return the code block",
			input:          "####hello world\n```go\npackage integrations\nhello\n",
			expectedOutput: "```go\npackage integrations\nhello\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := cleaner.FileWriteClean(tc.input)
			if output != tc.expectedOutput {
				t.Errorf("FileWriteClean returned incorrect output. Expected: %s, Got: %s", tc.expectedOutput, output)
			}
		})
	}
}
