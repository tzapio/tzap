package localwalker

import (
	"io"
	"io/fs"
	"os"

	"github.com/tzapio/tzap/pkg/types"
)

type LocalFile struct {
	filePath string
	file     *os.File
}

func NewLocalfile(filePath string) types.FileReader {
	return &LocalFile{filePath: filePath}
}

func (f *LocalFile) Filepath() string {
	return f.filePath
}
func (f *LocalFile) getFile() (*os.File, error) {
	if f.file == nil {
		file, err := os.Open(f.filePath)
		if err != nil {
			return nil, err
		}
		f.file = file
	}
	return f.file, nil
}
func (f *LocalFile) Open() (io.ReadCloser, error) {
	return f.getFile()
}
func (f *LocalFile) Close() (io.ReadCloser, error) {
	if f.file != nil {
		err := f.file.Close()
		if err != nil {
			return nil, err
		}
		f.file = nil
	}
	return f.getFile()
}
func (f *LocalFile) Stat() (fs.FileInfo, error) {
	file, err := f.getFile()
	if err != nil {
		return nil, err
	}
	return file.Stat()
}
