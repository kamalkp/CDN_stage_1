package main

import (
	"fmt"
	"os"
	"path/filepath"
	"log"
	"bufio"
	"strings"
	"github.com/ryanuber/go-glob"
	"flag"
)

var sitenamePtr = flag.String("sitename","examplecom","Port number to listen on")

//purge cache flags
var purgeMethodPtr = flag.String("purgemethod","","Any valid method")
var purgeFileNamePtr = flag.String("purgefilename","","Any valid filename")
var purgeFileExtPtr = flag.String("purgefileextension","","Any valid file extension starting with .")
var purgePathPtr = flag.String("purgepath","","Any valid location without regex or wildcards, starting with /")


func main() {
	searchDir := "/usr/local/nginx/cache/"

	flag.Parse()
	searchDir = searchDir + *sitenamePtr


	switch *purgeMethodPtr{

	case "entirecache" :

		fileList := []string{}
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)

			if !f.IsDir() {
				//			fmt.Println(path+f.Name())

				readFileByCachePath(path, "*")
			}

			return nil
		})

	case "filename" :
		fileList := []string{}
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)

			if !f.IsDir() {
				//			fmt.Println(path+f.Name())

				readFileByCachePath(path, "*"+*purgeFileNamePtr)
			}

			return nil
		})

	case "extension" :
		fileList := []string{}
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)

			if !f.IsDir() {
				//			fmt.Println(path+f.Name())

				readFileByCachePath(path, *purgeFileExtPtr)
			}

			return nil
		})

	case "path" :
		fileList := []string{}
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)

			if !f.IsDir() {
				//			fmt.Println(path+f.Name())

				readFileByCachePath(path, "*"+*purgePathPtr)
			}

			return nil
		})

	default :
		fmt.Println("error - Invalid purgemethod input")
		return
	}

/*	searchDir = searchDir + *sitenamePtr


	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)

		if !f.IsDir() {
//			fmt.Println(path+f.Name())

			readFileByCachePath(path, ".png")
		}

		return nil
	})*/

}



func readFileByCachePath(path string, search_str string){

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains( scanner.Text(), "KEY: " ){

			if glob.Glob(search_str, scanner.Text()) {
//				fmt.Println(path)
				file.Close()
				os.Remove(path)
			}

//			if strings.Contains(scanner.Text(), search_str){
//				fmt.Println(path)
//				file.Close()
//				os.Remove(path)
//				break
//			}
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}