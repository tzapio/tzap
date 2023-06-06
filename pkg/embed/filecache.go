package embed

import (
	"errors"
	"io"
	"io/fs"

	"github.com/tzapio/tzap/internal/logging/tl"
)

func (fe *Embedder) CheckFileCache() (changedFiles map[string]string, unchangedFiles map[string]int64, err error) {
	tl.Logger.Println("Checking file cache. Files:", len(fe.files))
	changedFiles = map[string]string{}
	unchangedFiles = map[string]int64{}

	for _, file := range fe.files {
		fileName := file.Filepath()

		fileStats, fileErr := file.Stat()
		if errors.Is(fileErr, fs.ErrNotExist) {
			tl.Logger.Println("File does not exist:", fileName)
			changedFiles[fileName] = ""
			continue
		}
		currentEditTime := fileStats.ModTime().UnixNano()
		cachedEditTime, exists := fe.filesTimestampsDB.Get(fileName)
		if exists && !isTimePassedSignificant(currentEditTime, cachedEditTime) {
			tl.Logger.Printf("NO CHANGE %s. Old Edittime: %d, New Edittime: %d, TimeDiff: %d", fileName, cachedEditTime, currentEditTime, cachedEditTime-currentEditTime)

			unchangedFiles[fileName] = cachedEditTime
			continue
		}
		readCloser, err := file.Open()
		if err != nil {
			println(err.Error())
			continue
		}
		fileContent, err := io.ReadAll(readCloser)
		if err != nil {
			println(err.Error())
			continue
		}
		fileContentStr := string(fileContent)
		tl.Logger.Printf("File %s has changed. Old Edittime: %d, New Edittime: %d, TimeDiff: %d", fileName, cachedEditTime, currentEditTime, cachedEditTime-currentEditTime)
		changedFiles[fileName] = fileContentStr
	}
	tl.Logger.Println("Finished checking file cache. Changed files:", len(changedFiles), "Unchanged files:", len(unchangedFiles))
	return changedFiles, unchangedFiles, nil
}
