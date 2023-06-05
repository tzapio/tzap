package cmdutil

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"time"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

// FileInZip represents a file in a zip archive.
type FileInZip struct {
	FilePath string
	Zipfile  *zip.File
}

func (f *FileInZip) Filepath() string {
	return f.FilePath
}

func (f *FileInZip) Open() (io.ReadCloser, error) {
	return f.Zipfile.Open()
}

func (f *FileInZip) Stat() (fs.FileInfo, error) {
	return virtualFileInfo{file: f.Zipfile}, nil
}

func (e *FileEvaluator) WalkDirFromURL(url string) ([]types.FileReader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	readerAt := bytes.NewReader(content)
	zipReader, err := zip.NewReader(readerAt, int64(len(content)))
	if err != nil {
		return nil, err
	}

	var list []types.FileReader
	for _, file := range zipReader.File {
		path := file.Name

		if !file.FileInfo().IsDir() && e.ShouldKeepPath(path) {
			tl.Logger.Println("KEEPFILE", path)
			list = append(list, &FileInZip{Zipfile: file, FilePath: path})
		} else {
			tl.Logger.Println("SKIPFILE", path)
		}
	}
	return list, nil
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
