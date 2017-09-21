package main

import (

	"strings"

	"fmt"
)

func main(){

	var test string = "sdvsd;sdcsd;addssdfs; proxy_cache_valid ..... ; sdfsdddddd; proxy_cache;--"

	substring(test, "proxy_cache_valid", "--")


}


func substring(input string, m string, n string){


	var start_index int
	var end_index int
	start_index = strings.Index(input, m)
//	fmt.Println(input[start_index:])
	end_index = strings.Index(input[start_index:], n)
	fmt.Println(end_index)

	fmt.Println(input[start_index:start_index+end_index+1])

}