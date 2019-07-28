package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Starting to watch your files 👀")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

}
