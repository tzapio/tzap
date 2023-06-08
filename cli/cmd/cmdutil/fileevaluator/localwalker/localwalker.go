package localwalker

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

type LocalWalker struct {
	dir string
	e   *fileevaluator.FileEvaluator
}

func New(e *fileevaluator.FileEvaluator, basePath string, dir string) *LocalWalker {
	return &LocalWalker{dir: dir, e: e}
}

func (f *LocalWalker) GetFiles() ([]types.FileReader, error) {
	var list []types.FileReader
	err := filepath.WalkDir(f.dir, func(path string, d fs.DirEntry, err error) error {
		tl.DeepLogger.Println("WALKDIR", f.dir, path, d.Name())
		if err != nil {
			return err
		}
		if d.IsDir() {
			if f.e.ShouldTraverseDir(path) || d.Name() == "." {
				tl.DeepLogger.Println("KEEPDIR", path)
				return nil
			} else {
				tl.DeepLogger.Println("SKIPDIR", path)
				return filepath.SkipDir
			}
		} else if f.e.ShouldKeepPath(path) {
			tl.Logger.Println("KEEPFILE", path)
			relPath, err := filepath.Rel(f.dir, path)
			if err != nil {
				return err
			}
			relPath = strings.TrimPrefix(relPath, "./")
			list = append(list, NewLocalfile(filepath.ToSlash(relPath)))
			return nil
		}
		tl.DeepLogger.Println("SKIPFILE", path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
