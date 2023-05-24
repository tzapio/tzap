package tl

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger
var DeepLogger *log.Logger
var UILogger *log.Logger
var UICompletionLogger *log.Logger

func init() {
	Logger = log.New(io.Discard, "Tzap ", log.Lshortfile)
	DeepLogger = log.New(io.Discard, "", 0)
	UILogger = log.New(io.Discard, "", 0)
	UICompletionLogger = log.New(io.Discard, "", 0)

}

func EnableLogger() {
	Logger.SetOutput(io.Writer(os.Stdout))

}
func EnableUILogger() {
	UICompletionLogger.SetOutput(io.Writer(os.Stdout))
}
func EnableUICompletionLogger() {
	UICompletionLogger.SetOutput(io.Writer(os.Stdout))
}

func EnableDeepLogger() {
	DeepLogger.SetOutput(io.Writer(os.Stdout))
}
