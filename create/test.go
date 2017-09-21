package main

import (
	"strings"
	"fmt"
	"strconv"
	"io/ioutil"
	"log"
)

func main() {

	input, err := ioutil.ReadFile("../manage/11175.conf")
	if err != nil {
		log.Fatalln(err)
	}



	lines := strings.Split(string(input), "\n")


	for j, line := range lines[0:]{

		words := strings.Split(line, " ")

		for i, word := range words{

			if word =="##pagespeed" {
				fmt.Println(strconv.Itoa(i) + " --word number " + word+" --line number "+strconv.Itoa(j))
			}
		}

	}




}
