package embed

import (
	"errors"
	"io"
	"io/fs"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

type FilestampCache struct {
	filesTimestampsDB types.DBCollectionInterface[int64]
}

func NewFilestampCache(filesTimestampsDB types.DBCollectionInterface[int64]) *FilestampCache {
	return &FilestampCache{filesTimestampsDB: filesTimestampsDB}
}
func (fc *FilestampCache) CheckFileCache(files []types.FileReader) (changedFiles map[string]string, unchangedFiles map[string]int64) {
	tl.Logger.Println("Checking file cache. Files:", len(files))
	changedFiles = map[string]string{}
	unchangedFiles = map[string]int64{}

	for _, file := range files {
		fileName := file.FilePath()

		fileStats, fileErr := file.Stat()
		if errors.Is(fileErr, fs.ErrNotExist) {
			tl.Logger.Println("File does not exist:", fileName)
			changedFiles[fileName] = ""
			continue
		}
		currentEditTime := fileStats.ModTime().UnixNano()
		cachedEditTime, exists := fc.filesTimestampsDB.Get(fileName)
		if exists && !isTimePassedSignificant(currentEditTime, cachedEditTime) {
			tl.DeepLogger.Printf("NO CHANGE %s. Old Edittime: %d, New Edittime: %d, TimeDiff: %d", fileName, cachedEditTime, currentEditTime, cachedEditTime-currentEditTime)
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
	return changedFiles, unchangedFiles
}

func (fc *FilestampCache) CacheFilestamps(embeddings *types.Embeddings, files []types.FileReader) error {
	if len(embeddings.Vectors) > 0 {
		var keyvals []types.KeyValue[int64]
		for _, vector := range embeddings.Vectors {
			for _, fileReader := range files {
				if fileReader.FilePath() == vector.Metadata.Filename {
					fileStat, err := fileReader.Stat()
					if err != nil {
						return err
					}
					keyvals = append(keyvals, types.KeyValue[int64]{Key: vector.Metadata.Filename, Value: fileStat.ModTime().UnixNano()})
				}
			}
		}
		added, err := fc.filesTimestampsDB.BatchSet(keyvals)
		if err != nil {
			panic("failing to store changed files should not happend and has probably caused some kind of corruption")
		}
		tl.Logger.Printf("Added %d files to file cache. Total: %d", added, len(embeddings.Vectors))
	}
	return nil
}
