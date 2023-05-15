package filelog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/tzapio/tzap/pkg/config"
)

type LogType string

const (
	TzapLog     LogType = "tzap"
	RequestLog  LogType = "request"
	ResponseLog LogType = "response"
)

func LogData(ctx context.Context, data interface{}, logType LogType) {
	config := config.FromContext(ctx)
	if config.LoggerOutput == "" {
		return
	}
	for i := 100000; i >= 0; i-- {
		dir := path.Join(config.LoggerOutput, string(logType))
		filename := path.Join(dir, fmt.Sprintf("out_%d.log", i))
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0700)
		}

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			writeJSONToFile(filename, data)
			fmt.Printf("Wrote log to file: %s\n", filename)
			return
		}
	}
	panic(fmt.Errorf("no free files?"))
}

func writeJSONToFile(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Errorf("error writing JSON to file: %w", err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		panic(fmt.Errorf("could not encode: %w", err))
	}

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("error reading file: %w", err))
	}

	if err := os.WriteFile(filename, fileBytes, 0644); err != nil {
		panic(fmt.Errorf("error writing file: %w", err))
	}

}
