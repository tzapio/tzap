package tzap

import (
	"fmt"
	"os"
	"strings"

	"github.com/tzapio/tzap/internal/logging/filelog"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/util"
)

// ChangeFilepath updates the filepath metadata in the Tzap data
func (t *Tzap) ChangeFilepath(filepath string) *Tzap {
	t.Data["filepath"] = filepath
	return t
}

// LoadCompletionOrRequestCompletionMD5 loads a task file if it exists and its MD5 checksum matches,
// otherwise requests a new task content from OpenAI and applies the changes to the original file
func (t *Tzap) LoadCompletionOrRequestCompletionMD5(filePath string) *Tzap {
	config := config.FromContext(t.C)
	md5memory := getMessageMD5(t)
	md5file := util.ReadFileP(filePath + ".md5")

	if _, err := os.Stat(filePath); err != nil {
		Log(t, "LoadTaskMD5 new file", filePath)
		return t.
			RequestChatCompletion().
			StoreCompletion(filePath)
	}
	if config.MD5Rewrites {
		if md5memory != md5file {
			include := false
			for _, rule := range config.MD5IncludeList {
				if strings.Contains(filePath, rule) {
					include = true
				}
			}
			if include {
				filelog.LogData(t.C, t, filelog.TzapLog)
				Log(t, "LoadTaskMD5 not matching", filePath)
				return t.
					RequestChatCompletion().
					StoreCompletion(filePath)
			}
		}
	}
	Log(t, "LoadTaskMD5 Matching", filePath, "enabled", config.MD5Rewrites)
	return t.LoadCompletion(filePath)

}

// getMessageMD5 generates an MD5 checksum of the messages content in a Tzap
func getMessageMD5(t *Tzap) string {
	tmpStr := ""
	for _, message := range GetMessages(t) {
		tmpStr += message.Content
	}
	return util.MD5Hash(tmpStr)
}

// writeMessageMD5 writes the MD5 checksum of messages content in a Tzap to a file
func writeMessageMD5(filename string, t *Tzap) error {
	md5memory := getMessageMD5(t)
	return os.WriteFile(filename+".md5", []byte(md5memory), 0755)
}
func CheckData(data types.MappedInterface) error {
	if _, ok := data["content"].(string); !ok {
		return fmt.Errorf("content key is not string: %+v", data["content"])
	}
	return nil
}
