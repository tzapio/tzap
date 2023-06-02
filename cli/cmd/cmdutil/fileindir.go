package cmdutil

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

// FileInDir represents a file in a directory.
type FileInDir struct {
	FilePath string
}

func (f *FileInDir) Filepath() string {
	return f.FilePath
}

func (f *FileInDir) Open() (io.ReadCloser, error) {
	return os.Open(f.FilePath)
}

func (e *FileEvaluator) ShouldTraverseDir(path string) bool {
	return !e.excludeMatcher.MatchesPath(path)
}

func (e *FileEvaluator) WalkDir(dir string) ([]types.FileReader, error) {
	var list []types.FileReader
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		tl.Logger.Println("WALKDIR", dir, path, d.Name())
		if err != nil {
			return err
		}
		if d.IsDir() {
			if e.ShouldTraverseDir(path) || d.Name() == "." {
				tl.Logger.Println("KEEPDIR", path)
				return nil
			} else {
				tl.Logger.Println("SKIPDIR", path)
				return filepath.SkipDir
			}
		} else if e.ShouldKeepPath(path) {
			tl.Logger.Println("KEEPFILE", path)
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(cwd, absPath)
			if err != nil {
				return err
			}
			relPath = strings.TrimPrefix(relPath, "./")
			list = append(list, &FileInDir{FilePath: filepath.ToSlash(relPath)})
			return nil
		}
		tl.Logger.Println("SKIPFILE", path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
