package zipwalker

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"

	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

// LocalWalker represents a file in a directory.
type ZipWalker struct {
	url              string
	relativeDirInZip string
	e                *fileevaluator.FileEvaluator
}

func New(e *fileevaluator.FileEvaluator, relativeDirInZip string, url string) *ZipWalker {
	return &ZipWalker{url: url, relativeDirInZip: relativeDirInZip, e: e}
}

func (z *ZipWalker) GetFiles() ([]types.FileReader, error) {
	tl.Logger.Println("Getting files: ", z.url)
	resp, err := http.Get(z.url)
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

		if !file.FileInfo().IsDir() && z.e.ShouldKeepPath(path) {
			tl.Logger.Println("KEEPFILE", path)
			list = append(list, &FileInZip{zipfile: file, filePath: path})
		} else {
			tl.Logger.Println("SKIPFILE", path)
		}
	}
	return list, nil
}
