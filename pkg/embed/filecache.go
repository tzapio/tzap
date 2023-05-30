package embed

import (
	"os"

	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/util"
)

func CheckFileCache(files []string) (map[string]string, map[string]string, error) {
	changedFiles := map[string]string{}
	unchangedFiles := map[string]string{}

	db, err := localdb.NewFileDB("./.tzap-data/filesmd5.db")
	if err != nil {
		return nil, nil, err
	}

	for _, fileName := range files {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			changedFiles[fileName] = ""
			continue
		}
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			continue
		}
		fileMD5 := util.MD5HashByte(fileContent)
		fileContentStr := string(fileContent)

		fileCache, exists := db.ScanGet(fileName)
		if exists && fileMD5 == fileCache.Value {
			unchangedFiles[fileName] = fileContentStr
			continue
		}

		changedFiles[fileName] = fileContentStr
	}
	println("celebrate")
	return changedFiles, unchangedFiles, nil
}
