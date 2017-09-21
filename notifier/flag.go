// [_Command-line flags_](http://en.wikipedia.org/wiki/Command-line_interface#Command-line_option)
// are a common way to specify options for command-line
// programs. For example, in `wc -l` the `-l` is a
// command-line flag.

package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import "flag"
import (
	"fmt"
	"strconv"
	"os"
	"io/ioutil"
	"log"
	"text/template"
)

var Path string
var DnsZoneFilePath string

var actionPtr = flag.String("action", "", "Action parameter indicating the action type/context")
var deviceIDPtr = flag.Int("device_id", 1111, "DeviceId value")
var domainNamePtr = flag.String("domainname", "example.com", "The domain name to be registered without the www")
var domainPortPtr = flag.Int("domainport",80,"Port number to listen on")
var sitenamePtr = flag.String("sitename","examplecom","Port number to listen on")
var ipAddressPtr = flag.String("ipaddress","127.0.0.1","Reverse proxy IP Address")
var clientIDPtr = flag.Int("client_id",1111,"This is the client ID")
var typePtr = flag.String("type","domain","Can be domain or service, i.e. domain level or service level")


func main() {

	// Basic flag declarations are available for string,
	// integer, and boolean options. Here we declare a
	// string flag `word` with a default value `"foo"`
	// and a short description. This `flag.String` function
	// returns a string pointer (not a string value);
	// we'll see how to use this pointer below.

	// This declares `numb` and `fork` flags, using a
	// similar approach to the `word` flag.
//	numbPtr := flag.Int("numb", 42, "an int")
//	boolPtr := flag.Bool("fork", false, "a bool")

	// It's also possible to declare an option that uses an
	// existing var declared elsewhere in the program.
	// Note that we need to pass in a pointer to the flag
	// declaration function.
//	var svar string
//	flag.StringVar(&svar, "svar", "bar", "a string var")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	// Here we'll just dump out the parsed options and
	// any trailing positional arguments. Note that we
	// need to dereference the pointers with e.g. `*wordPtr`
	// to get the actual option values.
	fmt.Println("Action Param:", *actionPtr)
	fmt.Println("Device ID:", *deviceIDPtr)
	fmt.Println("Domain Name:", *domainNamePtr)
	fmt.Println("Domain Port:", *domainPortPtr)
	fmt.Println("IP Address:", *ipAddressPtr)
	fmt.Println("Client ID:", *clientIDPtr)
	fmt.Println("Type:", *typePtr)
//	fmt.Println("numb:", *numbPtr)
//	fmt.Println("fork:", *boolPtr)
//	fmt.Println("svar:", svar)
//	fmt.Println("tail:", flag.Args())

	Path = strconv.Itoa(*deviceIDPtr)+".conf"
	DnsZoneFilePath = "dns_zone_file/pokecdn.net.json"

	//Define Flow based on action parameter
	if *actionPtr == "create"{
		fmt.Println("create request received...")
		createCDN()
		writeToDNSZoneFile()
	}

	if *actionPtr == "suspend"{
		fmt.Println("Suspend request received...")
		suspendCDN()
	}

	if *actionPtr == "unsuspend"{
		fmt.Println("Unsuspend request received...")
		unsuspendCDN()
	}

	if *actionPtr == "terminate"{
		fmt.Println("Terminate request received...")
		terminateCDN()
	}


	fmt.Println("success")

}





func createCDN(){
	//	fmt.Println(Env.Ipaddress)
	//	Path = strconv.Itoa(Env.Device_id)+".conf"

	// Define a template

	configFileByte, err := ioutil.ReadFile("template/template_for_create.conf") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	config_file_buffer := string(configFileByte) // convert content to a 'string'


	// prepare some data to insert into the template.
	type StructToFile struct {

		Client_id int
		Device_id int
		Domainname string
		Ipaddress string
		Domainport int
		Sitename string

	}

	var recipients = []StructToFile{
		{*clientIDPtr, *deviceIDPtr, *domainNamePtr, *ipAddressPtr, *domainPortPtr, *sitenamePtr},
		//              {"", "", false},
		//              {"", "", false},
	}

	// create a new template and parse the letter into it.
	t := template.Must(template.New("template").Parse(config_file_buffer))

	//        Path = strconv.Itoa(Env.Device_id)+".conf"

	createFile()
	//write to file - move to a separate function
	file, err := os.OpenFile(Path, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute( file , r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}


	//	createFile()
	//	writeFile()
	//	readFile()
	//	deleteFile()

}



func writeToDNSZoneFile(){
	//A very raw way of configuring string input to file, need to make this better
	insert_string := `,
"`+strconv.Itoa(*deviceIDPtr)+"."+*domainNamePtr+`": { "alias": "regular" }
  }
}
`


	f, err := os.OpenFile(DnsZoneFilePath, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	newPos, err := f.Seek(-7, 2); if err != nil {
		panic(err)
	}

	if _, err := f.WriteAt([]byte(insert_string), newPos); err != nil {
		panic(err)
	}

}



func suspendCDN(){

	err :=  os.Rename("/home/GoJSON/validation/"+Path, "/home/GoJSON/validation/suspend/"+Path+".suspend")

	if err != nil {
		fmt.Println(err)
		return
	}

}



func unsuspendCDN(){

	err :=  os.Rename("/home/GoJSON/validation/suspend/"+Path+".suspend", "/home/GoJSON/validation/"+Path)

	if err != nil {
		fmt.Println(err)
		return
	}

}



func terminateCDN(){

	err :=  os.Rename("/home/GoJSON/validation/"+Path, "/home/GoJSON/validation/terminate/"+Path+".terminate")

	if err != nil {
		fmt.Println(err)
		return
	}


}


func createFile() {
	// detect if file exists
	var _, err = os.Stat(Path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(Path)
		checkError(err) //okay to call os.exit()
		defer file.Close()
	}
}


func writeFile() {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(Path, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// write some text to file
	_, err = file.WriteString("This is a test\n")
	if err != nil {
		fmt.Println(err.Error())
		return //must return here for defer statements to be called
	}
	_, err = file.WriteString("mari belajar golang\n")
	if err != nil {
		fmt.Println(err.Error())
		return //same as above
	}

	// save changes
	err = file.Sync()
	if err != nil {
		fmt.Println(err.Error())
		return //same as above
	}
}

func readFile() {
	// re-open file
	var file, err = os.OpenFile(Path, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// read file
	var text = make([]byte, 1024)
	n, err := file.Read(text)
	if n > 0 {
		fmt.Println(string(text))
	}
	//if there is an error while reading
	//just print however much was read if any
	//at return file will be closed
}

func deleteFile() {
	// delete file
	var err = os.Remove(Path)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
