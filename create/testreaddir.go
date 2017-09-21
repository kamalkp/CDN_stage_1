package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	searchDir := "C:\\Users\\Getit\\Desktop\\cache2"

	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)

		if !f.IsDir() {
			fmt.Println(path+f.Name())
		}

		return nil
	})

}