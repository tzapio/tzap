package zipwalker

import (
	"archive/zip"
	"io"
	"io/fs"
	"time"
)

// FileInZip represents a file in a zip archive.
type FileInZip struct {
	filePath string
	zipfile  *zip.File
}

func (f *FileInZip) Filepath() string {
	return f.filePath
}

func (f *FileInZip) Open() (io.ReadCloser, error) {
	return f.zipfile.Open()
}

func (f *FileInZip) Stat() (fs.FileInfo, error) {
	return virtualFileInfo{file: f.zipfile}, nil
}

type virtualFileInfo struct {
	file *zip.File
}

func (v virtualFileInfo) Name() string {
	return v.file.Name
}

func (v virtualFileInfo) Size() int64 {
	return int64(v.file.UncompressedSize64)
}

func (v virtualFileInfo) Mode() fs.FileMode {
	return v.file.Mode()
}

func (v virtualFileInfo) ModTime() time.Time {
	return v.file.ModTime()
}

func (v virtualFileInfo) IsDir() bool {
	return v.file.FileInfo().IsDir()
}

func (v virtualFileInfo) Sys() interface{} {
	return nil
}
