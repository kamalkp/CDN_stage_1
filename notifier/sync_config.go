package main

import (
	"github.com/dietsche/rfsnotify"
	"os/exec"
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"strings"
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
					exec.Command("/root/CDN/copy.sh").Output()
					writeIPToNodeFiles()
					exec.Command("/root/CDN/copy_bck_sites.sh").Output()
					exec.Command("/root/CDN/SSL.sh").Output()

				}

			}
		}
	}()

	watcher.AddRecursive("/root/test/replicate/test/")
	//      if err != nil {
	//              log.Fatal(err)
	//      }
	<-done

}



func writeIPToNodeFiles(){

	commits := map[string]string{
		"/root/CDN/india/": "103.78.121.2",
		"/root/CDN/brazil/":   "181.41.201.175",
		"/root/CDN/australia/": "45.76.114.55",
		"/root/CDN/japan/": "45.76.48.30",
		"/root/CDN/singapore/": "45.77.34.118",
		"/root/CDN/southafrica/": "172.107.92.10",
		"/root/CDN/vancouver/": "68.168.30.2",
		"/root/CDN/eu_usa/": "172.106.22.7",
	}

	for k, v := range commits{
//		fmt.Printf("key[%s] value[%s]\n", k, v)
		changeListenerIP(k, v)

	}


	commits_presites := map[string][]string{
		"/root/CDN/back_sites/india/": {"103.78.121.4", "103.78.121.2", "103.78.121.3", "103.78.121.69" },
		"/root/CDN/back_sites/brazil/":   {"103.78.120.4", "181.41.217.232", "181.41.219.184", "181.41.201.175" },
		"/root/CDN/back_sites/australia/": {"103.78.120.4", "103.78.120.2", "103.78.120.3", "103.78.120.69" },
		"/root/CDN/back_sites/japan/": {"103.78.120.4", "103.78.120.2", "103.78.120.3", "103.78.120.69" },
		"/root/CDN/back_sites/singapore/": {"103.78.120.4", "103.78.120.2", "103.78.120.3", "103.78.120.69" },
		"/root/CDN/back_sites/southafrica/": {"165.73.240.130", "165.73.240.131", "165.73.240.132", "165.73.240.135" },
		"/root/CDN/back_sites/vancouver/": {"172.106.22.4", "172.106.4.19", "172.106.4.20", "172.106.22.69" },
		"/root/CDN/back_sites/eu_usa/": {"172.106.22.4", "172.106.4.19", "172.106.4.20", "172.106.22.69" },
	}

	for k, v := range commits_presites{
		//		fmt.Printf("key[%s] value[%s]\n", k, v)
		changeBckSitesIP(k, v)

	}


}



func changeListenerIP(NodePath string, ip string){

	files, _ := ioutil.ReadDir(NodePath)
	for _, filePath := range files {


		if !filePath.IsDir(){

		input, err := ioutil.ReadFile(NodePath+filePath.Name())
		if err != nil {
			log.Fatalln(err)
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {

			if strings.Contains(line, "172.106.22.7"){

				lines[i] = strings.Replace(line, "172.106.22.7", ip, -1)

			}

		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(NodePath+filePath.Name(), []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}

	}
	}

}



func changeBckSitesIP(NodePath string, ip []string){

	files, _ := ioutil.ReadDir(NodePath)
	for _, filePath := range files {


		if !filePath.IsDir(){

			input, err := ioutil.ReadFile(NodePath+filePath.Name())
			if err != nil {
				log.Fatalln(err)
			}

			lines := strings.Split(string(input), "\n")

			for i, line := range lines {

				if filePath.Name()=="psychz.net.conf"{

					if strings.Contains(line, "172.106.22.4"){

						lines[i] = strings.Replace(line, "172.106.22.4", ip[0], -1)

					}

				}

				if filePath.Name()=="photon.conf"{

					if strings.Contains(line, "172.106.4.19"){

						lines[i] = strings.Replace(line, "172.106.4.19", ip[1], -1)

					}

				}

				if filePath.Name()=="yard.conf"{

					if strings.Contains(line, "172.106.4.20"){

						lines[i] = strings.Replace(line, "172.106.4.20", ip[2], -1)

					}

				}

				if filePath.Name()=="dosattack.net.conf"{

					if strings.Contains(line, "172.106.22.69"){

						lines[i] = strings.Replace(line, "172.106.22.69", ip[3], -1)

					}

				}

			}
			output := strings.Join(lines, "\n")
			err = ioutil.WriteFile(NodePath+filePath.Name(), []byte(output), 0644)
			if err != nil {
				log.Fatalln(err)
			}

		}
	}

}