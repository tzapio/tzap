package types

import (
	"io"
)

// FileReader is an interface that wraps file reading operations.
type FileReader interface {
	// Open opens the file and returns a ReadCloser for accessing its contents.
	Open() (io.ReadCloser, error)
	// Name returns the name of the file.
	Filepath() string
}
