package cmdutil

import (
	"io"
	"log"
	"os"
)

var CMDUILogger *log.Logger

func init() {
	CMDUILogger = log.New(io.Discard, "", 0)
}

func EnableLogger() {
	CMDUILogger.SetOutput(io.Writer(os.Stdout))
}
