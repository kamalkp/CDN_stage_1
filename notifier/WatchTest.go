// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	//        "log"
	//        "fmt"
	"os/exec"
	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, _ := fsnotify.NewWatcher()
	//      if err != nil {
	//              log.Fatal(err)
	//      }
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//                              log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//                                      log.Println("modified file:", event.Name)
					exec.Command("/root/file.sh").Output()
				}
				//                      case err := <-watcher.Errors:
				//                              log.Println("error:", err)
			}
		}
	}()

	watcher.Add("test/")
	//      if err != nil {
	//              log.Fatal(err)
	//      }
	<-done
}
