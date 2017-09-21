package main

import (
	"github.com/dietsche/rfsnotify"
	"os/exec"
	"fmt"
	"os"
)


func main() {
	watcher, _ := rfsnotify.NewWatcher()
	//      if err != nil {
	//              log.Fatal(err)
	//      }
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
				case event := <-watcher.Events:
					if !true{
						fmt.Fprintln(os.Stdout, event)
					}
					if true{
						exec.Command("dir").Output()
					}

			}
		}
	}()

	watcher.AddRecursive("test/")
	//      if err != nil {
	//              log.Fatal(err)
	//      }
	<-done
}
