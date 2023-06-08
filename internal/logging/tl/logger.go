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
	Logger = log.New(io.Discard, "Tzap ", log.Lshortfile|log.Lmicroseconds)
	DeepLogger = log.New(io.Discard, "", 0)
	UILogger = log.New(io.Discard, "", 0)
	UICompletionLogger = log.New(io.Discard, "", 0)

}

func EnableLogger() {
	Logger.SetOutput(io.Writer(os.Stderr))
}
func EnableUILogger() {
	UICompletionLogger.SetOutput(io.Writer(os.Stderr))
}
func EnableUICompletionLogger() {
	UICompletionLogger.SetOutput(io.Writer(os.Stderr))
}
func EnableDeepLogger() {
	DeepLogger.SetOutput(io.Writer(os.Stderr))
}
