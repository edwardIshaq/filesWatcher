package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

var (
	watchFile  = flag.String("f", "./filesMap.json", "a json file with the list of files")
	sourceFile = flag.String("src", "", "the source file path")
	destFile   = flag.String("dst", "", "the destination file path used only with `src`")
)

func newFwatchError(err error, format string, a ...interface{}) fwatchError {
	return fwatchError{
		error:  err,
		Reason: fmt.Sprintf(format, a...),
	}
}

func (f fwatchError) String() string {
	return fmt.Sprintf("%s | %v", f.Reason, f.error)
}

func main() {
	flag.Parse()
	watchMap := map[string]string{}

	if watchFile != nil {
		manifestFile, err := os.Open(*watchFile)
		if err != nil {
			log.Println(newFwatchError(err, "Error openning manifest file"))
			return
		}
		bs, err := ioutil.ReadAll(manifestFile)
		if err != nil {
			log.Println(newFwatchError(err, "Error reading manifest file"))
			return
		}
		err = json.Unmarshal(bs, &watchMap)
		if err != nil {
			log.Println(newFwatchError(err, "Error parsing manifest file"))
			return
		}
	}

	log.Println("Starting to watch files")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				log.Printf("fileName %s: Op %v", event.Name, event.Op)
				go copyFile(event.Name, watchMap[event.Name])

				// watch for errors
			case err := <-watcher.Errors:
				log.Println("ERROR", err)
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
