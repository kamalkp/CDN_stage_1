// [_Command-line arguments_](http://en.wikipedia.org/wiki/Command-line_interface#Arguments)
// are a common way to parameterize execution of programs.
// For example, `go run hello.go` uses `run` and
// `hello.go` arguments to the `go` program.

package main

import "os"
import "fmt"
import "encoding/json"
import "github.com/asaskevich/govalidator"
//import "net"
import "io/ioutil"
import "strconv"
import "text/template"
import "log"



type JsonData struct {

	Id int `valid:"-"`
	Client_id int `valid:"-"`
	Device_id int `valid:"-"`
	Domainname string `valid:"dns"`
	Ipaddress string `valid:"ip"`
	Domainport int `valid:"-"`
	Sitename string `valid:"-"`
	Status string `valid:"-"`
	Updated_at string `valid:"-"`
	EnableSsl int `valid:"-"`
	EnableGzip int `valid:"-"`
	SslKey string `valid:"-"`
	SslCrt string `valid:"-"`
	WcRedirect int `valid:"-"`
	Action string `valid:"-"`

}

var Env JsonData //Global Config - Global Config Variable Names Always in Caps
var Path string
var DnsZoneFilePath string


func main() {

	// `os.Args` provides access to raw command-line
	// arguments. Note that the first value in this slice
	// is the Path to the program, and `os.Args[1:]`
	// holds the arguments to the program.


	argsWithProg := `{"id":2,"client_id":15881,"device_id":11065,"domainname":"dosattack.net","ipaddress":"151.101.192.253","domainport":443,"sitename":"dosattacknet","status":"Created","updated_at":"4/4/2017  11:05:31 AM","enableSsl":0,"enableGzip":0,"sslKey":"dosattacknet.key","sslCrt":"dosattacknet.crt","wcRedirect":0,"action":"create"}`//os.Args[1]

	fmt.Println(argsWithProg)

	var respBytes = []byte(argsWithProg)

	//	var Env JsonData
	if err := json.Unmarshal(respBytes, &Env); err != nil {
		fmt.Println(err)
	}


	if ok, err := govalidator.ValidateStruct(Env); err != nil { //Need to write custom validations for certain fields to extend the validator library
		panic(err)
	} else {
		fmt.Printf("OK: %v\n", ok)
	}

	//ip := net.ParseIP(Env.Ipaddress)


	Path = strconv.Itoa(Env.Device_id)+".conf"
	DnsZoneFilePath = "dns_zone_file/pokecdn.net.json"


	//Define Flow based on action parameter
	if Env.Action == "create"{
		fmt.Println("create request received...")
		createCDN(&Env)
		writeToDNSZoneFile(&Env)
	}

	if Env.Action == "suspend"{
		fmt.Println("Suspend request received...")
		suspendCDN(&Env)
	}

	if Env.Action == "unsuspend"{
		fmt.Println("Unsuspend request received...")
		unsuspendCDN(&Env)
	}

	if Env.Action == "terminate"{
		fmt.Println("Terminate request received...")
		terminateCDN(&Env)
	}



	fmt.Println("success")

}//end main



func createCDN(env* JsonData){
	//	fmt.Println(Env.Ipaddress)
	//	Path = strconv.Itoa(Env.Device_id)+".conf"

	// Define a template

	configFileByte, err := ioutil.ReadFile("template/template_new.conf") // just pass the file name
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
		Status string
		Updated_at string
		EnableSsl int
		EnableGzip int
		SslKey string
		SslCrt string
		WcRedirect int

	}

	var recipients = []StructToFile{
		{env.Client_id, env.Device_id, env.Domainname, env.Ipaddress, env.Domainport, env.Sitename, env.Status, env.Updated_at, env.EnableSsl, env.EnableGzip, env.SslKey, env.SslCrt, env.WcRedirect},
		//              {"", "", false},
		//              {"", "", false},
	}

	// create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Parse(config_file_buffer))

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



func writeToDNSZoneFile(env* JsonData){
	//A very raw way of configuring string input to file, need to make this better
	insert_string := `,
"`+strconv.Itoa(env.Device_id)+"."+env.Domainname+`": { "alias": "regular" }
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



func suspendCDN(env* JsonData){

	err :=  os.Rename("/home/GoJSON/validation/"+Path, "/home/GoJSON/validation/suspend/"+Path+".suspend")

	if err != nil {
		fmt.Println(err)
		return
	}

}



func unsuspendCDN(env* JsonData){

	err :=  os.Rename("/home/GoJSON/validation/suspend/"+Path+".suspend", "/home/GoJSON/validation/"+Path)

	if err != nil {
		fmt.Println(err)
		return
	}

}



func terminateCDN(env* JsonData){

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
