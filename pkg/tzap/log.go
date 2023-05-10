package tzap

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/util"
)

func recentStack() string {
	// Get the caller information for the two most recent calls
	pc := make([]uintptr, 4)
	runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc)
	r := ""
	// Print the filename and line number for each call
	for i := 0; i < 4; i++ {
		frame, more := frames.Next()
		if !more {
			break
		}
		r = r + fmt.Sprintf(" %s:%d", frame.File, frame.Line)
	}
	return r
}

// Package-level variables for global throttling
var (
	LastExecutionTime time.Time
	globalMutex       sync.Mutex
	MessageBuffer     []string
)

func Log(t *Tzap, messages ...interface{}) {
	configuration := config.FromContext(t.C)
	if !configuration.EnableLogs {
		return
	}
	globalMutex.Lock()
	defer globalMutex.Unlock()

	now := time.Now()

	// Logging logic
	// Replace GetNames(t) and printWithSpace with their respective implementations
	names := GetNames(t)
	paths := strings.Join(names, "-")

	spacer := "---> " + fmt.Sprintf("%d", (len(names))) + util.CreateSpaces(len(names))
	message := fmt.Sprintf("%s%s\n%s\\____", spacer, paths, util.CreateSpaces(len(spacer+paths)))
	message += sprintWithSpace(messages...) + " " + recentStack()
	message += fmt.Sprintf("\n%s\\____\n", util.CreateSpaces(len(spacer+paths)+5))

	// if time has passed since last message print without adding.

	// Add the message to the buffer
	MessageBuffer = append(MessageBuffer, message)

	// If the buffer has more than 15 messages, remove the oldest one
	if len(MessageBuffer) > 15 {
		MessageBuffer = MessageBuffer[1:]
	}
	if now.Sub(LastExecutionTime) >= 1*time.Second {
		LastExecutionTime = now
		rawFlush()
	}

}

// no locks. Make sure you locked GlobalMutex.
func rawFlush() {
	for _, message := range MessageBuffer {
		fmt.Print(message)
	}
	// Clear the message buffer
	MessageBuffer = []string{}
}
func Flush() {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	LastExecutionTime = time.Now()
	rawFlush()
}
func ResetFlush() {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	LastExecutionTime = time.Time{}
	rawFlush()
}
func sprintWithSpace(a ...interface{}) string {
	t := ""
	for i, arg := range a {
		t += fmt.Sprintf("%v", arg)
		if i < len(a)-1 {
			t += " "
		}
	}
	return t
}
func Logf(t *Tzap, format string, args ...any) {
	Log(t, fmt.Sprintf(format, args...))
}
func GetNames(t *Tzap) []string {
	var names []string

	if t.Parent != nil {
		names = GetNames(t.Parent)
	}

	if t.Name == "" {
		panic("Name required!")
	}

	return append(names, t.Name)

}
