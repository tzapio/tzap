package embed

import (
	"os"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/util"
)

func CheckFileCache(files []string) (changedFiles map[string]string, unchangedFiles map[string]string, err error) {
	changedFiles = map[string]string{}
	unchangedFiles = map[string]string{}
	db, err := localdb.NewFileDB("./.tzap-data/filesmd5.db")
	if err != nil {
		return
	}

	for _, fileName := range files {
		if _, fileErr := os.Stat(fileName); os.IsNotExist(fileErr) {
			changedFiles[fileName] = ""
			continue
		}
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			continue
		}
		currentMD5 := util.MD5HashByte(fileContent)
		fileContentStr := string(fileContent)

		cachedMD5, exists := db.Get(fileName)
		if exists && currentMD5 == cachedMD5 {
			unchangedFiles[fileName] = fileContentStr
			continue
		}
		tl.Logger.Printf("File %s has changed. Old MD5: %s, New MD5: %s", fileName, cachedMD5, currentMD5)
		changedFiles[fileName] = fileContentStr
	}

	return
}
