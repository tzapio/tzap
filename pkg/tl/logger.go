package tl

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	Logger = log.New(io.Discard, "Tzap ", log.Lshortfile)
}

func EnableLogger() {
	Logger.SetOutput(io.Writer(os.Stdout))
}
