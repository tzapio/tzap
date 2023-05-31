package embed

import (
	"os"

	"github.com/tzapio/tzap/internal/logging/tl"
)

func (fe *Embedder) CheckFileCache(files []string) (changedFiles map[string]string, unchangedFiles map[string]int64, err error) {
	changedFiles = map[string]string{}
	unchangedFiles = map[string]int64{}

	for _, fileName := range files {
		fileStats, fileErr := os.Stat(fileName)
		if os.IsNotExist(fileErr) {
			changedFiles[fileName] = ""
			continue
		}
		currentEditTime := fileStats.ModTime().UnixNano()

		cachedEditTime, exists := fe.filesTimestampsDB.Get(fileName)
		if exists && !isTimePassedSignificant(currentEditTime, cachedEditTime) {
			unchangedFiles[fileName] = cachedEditTime
			continue
		}
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			continue
		}
		fileContentStr := string(fileContent)
		tl.Logger.Printf("File %s has changed. Old Edittime: %d, New Edittime: %d, TimeDiff: %d", fileName, cachedEditTime, currentEditTime, cachedEditTime-currentEditTime)
		changedFiles[fileName] = fileContentStr
	}
	return changedFiles, unchangedFiles, nil
}
