package types

import (
	"io"
	"io/fs"
)

// FileReader is an interface that wraps file reading operations.
type FileReader interface {
	// Open opens the file and returns a ReadCloser for accessing its contents.
	Open() (io.ReadCloser, error)
	// Name returns the name of the file.
	Filepath() string
	Stat() (fs.FileInfo, error)
}
type FileWalker interface {
	GetFiles() ([]FileReader, error)
}
