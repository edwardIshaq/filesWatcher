package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
)

type (
	fwatchError struct {
		error
		Reason string
	}

	fileMap struct {
		SourcePath string
		DestPath   string
	}
)

func newFwatchError(err error, format string, a ...interface{}) fwatchError {
	return fwatchError{
		error:  err,
		Reason: fmt.Sprintf(format, a...),
	}
}

func main() {
	fmt.Println("Starting to watch your files 👀")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	watchMap := map[string]string{
		"/tmp/watchme": "/tmp/backup/watchme",
	}

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("fileName %s: Op %v", event.Name, event.Op)
				go copyFile(event.Name, watchMap[event.Name])

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	for src := range watchMap {
		// out of the box fsnotify can watch a single file, or a single directory
		if err := watcher.Add(src); err != nil {
			fmt.Println("ERROR", err)
		}
	}
	<-done
}

func copyFile(src, dst string) error {
	orig, err := os.Open(src)
	defer orig.Close()
	if err != nil {
		return newFwatchError(err, "Error opening original file, maybe it was deleted")
	}
	origStat, err := orig.Stat()
	if err != nil {
		return newFwatchError(err, "Error Reading orig.Stat()")
	}

	copyDest, err := os.Create(dst)
	if err != nil {
		return newFwatchError(err, "Error Opening Dest file")
	}

	_, err = io.Copy(copyDest, orig)
	if err != nil {
		return newFwatchError(err, "Issue copying file")
	}

	err = copyDest.Chmod(origStat.Mode())
	if err != nil {
		return newFwatchError(err, "Error while chmod the dest file")
	}
	return nil
}
