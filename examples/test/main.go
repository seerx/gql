package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func main() { 
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("event: ",  event)
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/path/to/file1")
	if err != nil {
		panic(err)
	}
	err = watcher.Add("/Users/dotjava/workspace/go-projects")//也可以监听文件夹
	if err != nil {
		panic(err)
	}
	<-done
}
