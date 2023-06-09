package cmdui

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tzapio/tzap/internal/logging/tl"
)

func (ui *CMDUI) WatchSavesToFile(changedWithin time.Duration, times int) error {

	// Watch for file changes
	var lastChangeTime time.Time
	var changes int

	// Create a new file watcher to monitor changes to the file
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Start watching the file
	err = watcher.Add(ui.filePath)
	if err != nil {
		return err
	}
	defer watcher.Remove(ui.filePath)

	for {
		select {
		case event := <-watcher.Events:
			// If the file was modified, check how long ago the last modification was
			if event.Op.Has(fsnotify.Write) {
				now := time.Now()
				passed := now.Sub(lastChangeTime)
				if passed < time.Millisecond*50 {
					continue
				}
				if passed <= changedWithin {
					changes++
					tl.Logger.Println("file saved - incrementing changes", changes)
				} else {
					tl.Logger.Println("file saved - resetting changes")
					changes = 1
				}
				lastChangeTime = now

				// If we've seen multiple changes within the time window, call the callback
				if changes >= times {
					tl.Logger.Printf("saved %d times in 2 seconds", times)
					changes = 0
					return nil
				}
			}
		case err := <-watcher.Errors:
			println("error:", err)
			return err
		}
	}
}
