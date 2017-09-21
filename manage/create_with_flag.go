
package main

import "flag"
import (
	"fmt"
	"strconv"
	"os"
	"io/ioutil"
	"log"
	"text/template"
	"strings"

	"bufio"
	"path/filepath"
	"os/exec"
)
//Paths to files and actions
var TemplatePath string = ""
var TemplateFilePath string = "../template/template_for_create.conf"
var SuspendPath string = "suspend/"
var TerminatePath string = "terminate/"
var DnsZoneFilePath string = "../dns_zone_file/pokecdn.net.json"
var CacheIncludeFilePath string = "../cache/cache.conf"
var Path_To_Cache_Directory = "/usr/local/nginx/cache/"

//create flags - create domain, subdomain
var actionPtr = flag.String("action", "", "Action parameter indicating the action type/context")
var deviceIDPtr = flag.Int("device_id", 0, "DeviceId value")
var domainNamePtr = flag.String("domainname", "example.com", "The domain name to be registered without the www")
var domainPortPtr = flag.Int("domainport",80,"Port number to listen on")
var sitenamePtr = flag.String("sitename","examplecom","Port number to listen on")
var ipAddressPtr = flag.String("ipaddress","127.0.0.1","Reverse proxy IP Address")
var clientIDPtr = flag.Int("client_id",0,"This is the client ID")
var typePtr = flag.String("type","domain","Can be domain or service, i.e. domain level or service level")

//create flags - domain level

//Routes
var route = flag.String("route","","Routes can be premium or regular")

//botnet and hotlink flags
var botnetPtr = flag.String("botnet","","Can be true or false, i.e. domain level or service level")
var hotlinkPtr = flag.String("hotlink","","Can be true or false, i.e. domain level or service level")


//Gzip
var gzipPtr = flag.String("gzip","","Can be true or false, i.e. domain level or service level")
var gzipCustomMimeTypePtr = flag.String("gzipmimetype","","Can be anything, i.e. domain level only")
var gzipRemoveMimePtr = flag.String("gzipmimeremove","","Can be anything, i.e. domain level only")
var gzipLocationPtr = flag.String("gziplocation","","starts with /somelocation - with no spaces after /")
var gzipExcludeMimePtr = flag.String("gzipexcludemime","","starts with /somelocation - with no spaces after /")


//pagespeed
var pagespeedPtr = flag.String("pagespeed","","can be true or false")
var pagespeedLocationPtr = flag.String("pagespeedlocation","","must be a location name without any slashes or asterisks")
var removeByLocationPtr = flag.String("removebylocation","","true or false")


//access control flags
var accessControlLocationPtr = flag.String("accesscontrollocation","","must be a location name with a preceding / and no space")
var accessControlWhitelistIPPtr = flag.String("whitelistip","","must be a valid IP address")
var accessControlBlacklistIPPtr = flag.String("blacklistip","","must be a valid IP address")
var removeBlacklistIPPtr = flag.String("removeblacklist","","must be true if used")
var removeWhitelistIPPtr = flag.String("removewhitelist","","must be true if used")


//Manage domain flags - subdomain, parked domain, bandwidth
var subdomainPtr = flag.String("subdomain","","should be a valid subdomain value")
var subdomainPortPtr = flag.Int("subdomainport",80,"Port number to listen on")
var subdomainipAddressPtr = flag.String("subdomainipaddress","127.0.0.1","Reverse proxy IP Address")
var subdomainTypePtr = flag.String("subdomaintype","","can be alias or origin")
var removesubdomainPtr = flag.String("removesubdomain","","can be alias or origin")
var subdomainsitenamePtr = flag.String("subdomainsitename","","subdomainsitename")


//parked domain
var parkeddomainPtr = flag.String("parkeddomain","","should be a mirror")
var removeparkeddomainPtr = flag.String("removeparkeddomain","","should be a mirror")
var bandwidthPtr = flag.Int("bandwidth",-1,"Usage in megabytes, must be an integer")
var bandwidthUnitPtr = flag.String("bandwidthunit","","Unit in megabytes, gbs must be m , k")


//cache flags
var cacheTimePtr = flag.String("cachetime","","Unit must be an integer")
var cacheTimeUnitPtr = flag.String("cachetimeunit","","Units must be ms, s, m, h, d, w, M, y")
var cacheBypassPtr = flag.String("cachebypasslocation","","Any valid location without regex or wildcards, starting with /")
var removeCacheBypassPtr = flag.String("removecachebypass","","true or false")
var cacheTTLTimePtr = flag.String("cachettltime","","Unit must be an integer")
var cacheTTLTimeUnitPtr = flag.String("cachettltimeunit","","Units must be ms, s, m, h, d, w, M, y")
var cacheTTLLocationPtr = flag.String("cachettllocation","","Any valid location without regex or wildcards, starting with /")
var removeTTLPtr = flag.String("removettl","","true or false")
var staleContentPtr = flag.String("stalecontentlocation","","Any valid location without regex or wildcards, starting with /")
var removeStaleContentPtr = flag.String("removestalecontent","","true or false")
var stringQueryPtr = flag.String("stringquerylocation","","Any valid location without regex or wildcards, starting with /")
var removeStringQueryPtr = flag.String("removestringquery","","true or false")
var remove_cache_location = flag.String("remove_cache_location","","location value")
var enable_global_cachePtr = flag.String("enable_global_cache","","true or false")

var sslOnPtr = flag.String("sslon","","true or false")


//purge cache flags
var purgeMethodPtr = flag.String("purgemethod","","Any valid method")
var purgeFileNamePtr = flag.String("purgefilename","","Any valid filename")
var purgeFileExtPtr = flag.String("purgefileextension","","Any valid file extension starting with .")
var purgePathPtr = flag.String("purgepath","","Any valid location without regex or wildcards, starting with /")


//ssl
var sslEnabledPtr = flag.String("sslenable","","can be true or false")
var sslTypePtr = flag.String("ssltype","","can be auto or thirdparty")
var privateKeyPtr = flag.String("privatekey","","Any valid location without regex or wildcards, starting with /")
var certPtr = flag.String("cert","","Any valid location without regex or wildcards, starting with /")


//layer7
var layer7Ptr = flag.String("layer7","","can be true or false")
var actionNamePtr = flag.String("action_name","","can be layer7")
var enablePtr = flag.String("enable","","can be true or false")
var layer7LocationPtr = flag.String("location","","layer 7 Exempt Location")
var ip_whitelistPtr = flag.String("ip_whitelist","","layer 7 whitelist IP Address")
var layer7whitelistPtr = flag.String("layer7whitelist","","layer 7 whitelist enable/disable")


func main() {

	flag.Parse()
//	fmt.Println("The error was "+ string(err))

	TemplatePath = TemplatePath+strconv.Itoa(*deviceIDPtr)+".conf"
	SuspendPath = SuspendPath+strconv.Itoa(*deviceIDPtr)+".conf"
	TerminatePath = TerminatePath+strconv.Itoa(*deviceIDPtr)+".conf"

	//Define Flow based on action parameter
	if *actionPtr == "create"{
		fmt.Println("create request received...")
		if checkIfFileExists(){
			addToCDN()
		} else{
			createCDN()
		}
		writeToDNSZoneFile()
	}

	if *actionPtr == "suspend" && *typePtr=="service"{
		fmt.Println("Suspend request received...")
		suspendCDN()
	}

	if *actionPtr == "unsuspend" && *typePtr=="service"{
		fmt.Println("Unsuspend request received...")
		unsuspendCDN()
	}

	if *actionPtr == "terminate" && *typePtr=="service"{
		fmt.Println("Terminate request received...")
		terminateCDN()
	}

	if *actionPtr == "suspend" && *typePtr=="domain"{
		fmt.Println("Domain Suspend request received...")
		domainLevelSuspendCDN()
	}

	if *actionPtr == "unsuspend" && *typePtr=="domain"{
		fmt.Println("Domain Unsuspend request received...")
		domainLevelUnsuspendCDN()
	}

	if *actionPtr == "terminate" && *typePtr=="domain"{
		fmt.Println("Domain Terminate request received...")
		domainLevelTerminateCDN()
	}

	if *route != ""{
		fmt.Println("Change Route request received...")
		changeDNSRoute()
	}

	if *botnetPtr == "true"{
		fmt.Println("Activate botnet request received...")
		activateBotnet()
	}

	if *botnetPtr == "false"{
		fmt.Println("Deactivate botnet request received...")
		deactivateBotnet()
	}

	if *hotlinkPtr == "true"{
		fmt.Println("Activate hotlink request received...")
		activateHotlink()
	}

	if *hotlinkPtr == "false"{
		fmt.Println("Deactivate hotlink request received...")
		deactivateHotlink()
	}

	if *gzipPtr == "true"{
		fmt.Println("Gzip request received...")
		enableGzip()
	}

	if *gzipPtr == "false"{
		fmt.Println("Disable Gzip request received...")
		disableGzip()
	}

	if *gzipCustomMimeTypePtr != "" && *gzipRemoveMimePtr=="" && *gzipExcludeMimePtr == "" && *gzipLocationPtr == ""{
		fmt.Println("Gzip custom mime type request received...")
		addGzipMimeType()
	}

	if *gzipCustomMimeTypePtr != "" && *gzipRemoveMimePtr=="true"{
		fmt.Println("Remove Gzip custom mime type request received...")
		removeGzipMimeType()
	}

	if *gzipLocationPtr != "" && *gzipExcludeMimePtr == ""{
		fmt.Println("Add Gzip Custom location Mimetype request received...")
		gzipLocationMimeExclusion(*gzipLocationPtr, *gzipCustomMimeTypePtr)
	}

	if *gzipLocationPtr != "" && *gzipExcludeMimePtr == "true"{
		fmt.Println("Remove Gzip Custom location Mimetype request received...")
		removeGzipLocationMimeExclusion(*gzipLocationPtr, *gzipCustomMimeTypePtr)
	}

	if *pagespeedPtr == "true" && *pagespeedLocationPtr == ""{
		fmt.Println("Pagespeed request received...")
		createPageSpeed()
	}

	if *pagespeedPtr == "false" && *pagespeedLocationPtr == ""{
		fmt.Println("Pagespeed request received...")
		removePageSpeed()
	}

	if *pagespeedLocationPtr != "" && *removeByLocationPtr == ""{
		fmt.Println("Pagespeed disallow by location request received...")
		blockPageSpeedByLocation()
	}

	if *pagespeedLocationPtr != "" && *removeByLocationPtr == "true"{
		fmt.Println("Pagespeed remove disallow by location request received...")
		unblockPageSpeedByLocation()
	}

	if *subdomainPtr != ""{
		fmt.Println("Subdomain request received...")
		if *subdomainTypePtr == "origin" {
			createSubDomain()
		}
		if *subdomainTypePtr == "alias" && *removesubdomainPtr=="" {
			fmt.Println("Create Alias Subdomain request received...")
			createSubDomainAlias()
		}
		if *subdomainTypePtr == "alias" && *removesubdomainPtr=="true"{
			fmt.Println("Remove Alias Subdomain request received...")
			removeSubDomainAlias()
		}

	}

	if *parkeddomainPtr != "" && *removeparkeddomainPtr==""{
		fmt.Println("Park domain request received...")
		configParkedDomain()
	}

	if *removeparkeddomainPtr == "true"{
		fmt.Println("Remove Park domain request received...")
		removeParkedDomain()
	}

	if *bandwidthPtr > 0{
		fmt.Println("Bandwidth request received...")
		configBandwidth()
	}

	if *bandwidthPtr == 0 && bandwidthPtr!=nil{
		fmt.Println("Remove Bandwidth limit request received...")
		removeBandwidthLimit()
	}

	if *accessControlLocationPtr != "" && *accessControlWhitelistIPPtr != "" && *removeBlacklistIPPtr == "" && *removeWhitelistIPPtr == ""{
		fmt.Println("Access Control whitelist request received...")
		addAccessControlAtLocation(*accessControlLocationPtr, *accessControlWhitelistIPPtr)
	}

	if *accessControlLocationPtr != "" && *accessControlBlacklistIPPtr != "" && *removeBlacklistIPPtr == "" && *removeWhitelistIPPtr == ""{
		fmt.Println("Access Control blacklist request received...")
		addAccessControlAtLocation(*accessControlLocationPtr, *accessControlBlacklistIPPtr)
	}

	if *removeBlacklistIPPtr != ""{
		fmt.Println("Remove blacklist request received...")
		removeBlackListAtLocation(*accessControlLocationPtr, *accessControlBlacklistIPPtr)
	}

	if *removeWhitelistIPPtr != ""{
		fmt.Println("Remove whitelist request received...")
		removeWhiteListAtLocation(*accessControlLocationPtr, *accessControlWhitelistIPPtr)
	}

	if *cacheTimePtr != "" && *cacheTimeUnitPtr != "" {
		fmt.Println("Cache set time request received...")
		makeCacheGlobalEntry()
	}

	if *cacheBypassPtr != "" && *removeCacheBypassPtr == ""{
		fmt.Println("Cache bypass request received...")
		bypassCacheForLocation(*cacheBypassPtr)
	}

	if *cacheBypassPtr != "" && *removeCacheBypassPtr == "true"{
		fmt.Println("Remove Cache bypass request received...")
		removeBypassCacheForLocation(*cacheBypassPtr)
	}

	if *cacheTTLTimePtr != "" && *cacheTTLTimeUnitPtr != "" && *removeTTLPtr == "" {
		fmt.Println("Cache set TTL time request received...")
		makeCacheTTLForLocation(*cacheTTLLocationPtr, *cacheTTLTimePtr, *cacheTTLTimeUnitPtr)
	}

	if *removeTTLPtr == "true" {
		fmt.Println("Cache remove TTL time request received...")
		removeCacheTTLForLocation(*cacheTTLLocationPtr, *cacheTTLTimePtr, *cacheTTLTimeUnitPtr)
	}

	if *staleContentPtr != "" && *removeStaleContentPtr == ""{
		fmt.Println("Add Stale Content request received...")
		enableStaleContentForLocation(*staleContentPtr)
	}

	if *staleContentPtr != "" && *removeStaleContentPtr == "true"{
		fmt.Println("Remove Stale Content request received...")
		disableStaleContentForLocation(*staleContentPtr)
	}

	if *stringQueryPtr != "" && *removeStringQueryPtr == ""{
		fmt.Println("Add String Query Cache request received...")
		enableStringQueryCacheForLocation(*stringQueryPtr)
	}

	if *stringQueryPtr != "" && *removeStringQueryPtr == "true"{
		fmt.Println("Remove String Query Cache request received...")
		disableStringQueryCacheForLocation(*stringQueryPtr)
	}

	if *remove_cache_location != ""{
		fmt.Println("Remove Cache for location request received...")
		removeCacheForLocation()
	}

	if *enable_global_cachePtr == "true" {
		fmt.Println("Enable Global Cache request received...")
		enableAllCacheOptions()
	}

	if *enable_global_cachePtr == "false" {
		fmt.Println("Disable Global Cache request received...")
		disableAllCacheOptions()
	}

	if *purgeMethodPtr != ""{
		fmt.Println("Purge Cache request received...")
		purgeCache()
	}

	if *sslEnabledPtr == "true" && *sslTypePtr=="auto" {
		fmt.Println("Enable auto SSL request received...")
		enableAutoSSL()
	}

	if *sslEnabledPtr == "false" && *sslTypePtr=="auto" {
		fmt.Println("Disable auto SSL request received...")
		disableAutoSSL()
	}

	if *sslEnabledPtr == "true" && *sslTypePtr=="thirdparty" {
		fmt.Println("Enable third party SSL request received...")
		enableThirdPartySSL()
	}

	if *sslEnabledPtr == "false" && *sslTypePtr=="thirdparty" {
		fmt.Println("Disable third party SSL request received...")
		disableThirdPartySSL()
	}

	if *layer7Ptr == "true" {
		fmt.Println("Enable Layer7 request received...")
		enableLayer7()
	}

	if *layer7Ptr == "false" {
		fmt.Println("Disable Layer7 request received...")
		disableLayer7()
	}

	if *actionNamePtr == "layer7" && *enablePtr == "" && *ip_whitelistPtr == "" {
		fmt.Println("Layer7 Exemption Path Enable request received...")
		layer7ExemptionPathEnable()
	}

	if *actionNamePtr == "layer7" && *enablePtr == "false" && *ip_whitelistPtr == "" {
		fmt.Println("Layer7 Exemption Path Disable request received...")
		layer7ExemptionPathDisable()
	}

	if *actionNamePtr == "layer7" && *ip_whitelistPtr != "" && *layer7whitelistPtr == "" {
		fmt.Println("Layer7 add whitelist request received...")
		layer7WhiteListAdd()
	}

	if *actionNamePtr == "layer7" && *ip_whitelistPtr != "" && *layer7whitelistPtr == "disable" {
		fmt.Println("Layer7 remove whitelist request received...")
		layer7WhiteListRemove()
	}


}



func createCDN(){
	//      fmt.Println(Env.Ipaddress)
	//      Path = "/root/test/CDN/"+strconv.Itoa(Env.Device_id)+".conf"

	// Define a template

	configFileByte, err := ioutil.ReadFile(TemplateFilePath) // just pass the file name
	if err != nil {
		fmt.Print("error : "+err.Error())
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
	file, err := os.OpenFile(TemplatePath, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute( file , r)
		if err != nil {
			fmt.Print("error : "+err.Error())
			log.Println("executing template:", err)
		}
	}


}



func addToCDN(){
//append to an existing deviceid file
//create a buffer file from the template - merge the 2 files
	configFileByte, err := ioutil.ReadFile(TemplateFilePath) // just pass the file name
	if err != nil {
		fmt.Print("error : "+err.Error())
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
	file, err := os.OpenFile(TemplatePath, os.O_RDWR|os.O_APPEND, 0660)
	checkError(err)
	defer file.Close()

	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute( file , r)
		if err != nil {
			fmt.Print("error : "+err.Error())
			log.Println("executing template:", err)
		}
	}


}



func createSubDomain(){
	//append to an existing deviceid file
	//create a buffer file from the template - merge the 2 files
	configFileByte, err := ioutil.ReadFile(TemplateFilePath) // just pass the file name
	if err != nil {
		fmt.Print("error : "+err.Error())
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
		{*clientIDPtr, *deviceIDPtr, *subdomainPtr, *subdomainipAddressPtr, *subdomainPortPtr, *subdomainsitenamePtr},
		//              {"", "", false},
		//              {"", "", false},
	}

	// create a new template and parse the letter into it.
	t := template.Must(template.New("template").Parse(config_file_buffer))

	//        Path = strconv.Itoa(Env.Device_id)+".conf"

	createFile()
	//write to file - move to a separate function
	file, err := os.OpenFile(TemplatePath, os.O_RDWR|os.O_APPEND, 0660)
	checkError(err)
	defer file.Close()

	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute( file , r)
		if err != nil {
			fmt.Print("error : "+err.Error())
			log.Println("executing template:", err)
		}
	}


}



func createSubDomainAlias(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		//This ensures that the parked domain is appended to the said domain only
		if strings.Contains(line, "##server_name") && strings.Contains(lines[i+1], *domainNamePtr){
			lines[i+1] = strings.Replace(lines[i+1], ";", " "+*subdomainPtr, -1) + " ;"
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeSubDomainAlias(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
//	var final_line string

	for i, line := range lines {
		//This ensures that the parked domain is appended to the said domain only
		if strings.Contains(line, "##"+*domainNamePtr){
			j, _ := lineNumForString(i,"##server_name", TemplatePath)
			//			fmt.Println(j)


			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), *subdomainPtr)//final_line//strings.Replace(lines[j+1], *parkeddomainPtr, "", -1)
			//			lines[i+1] = strings.Replace(lines[i+1], "", "", -1) + " ;"
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func writeToDNSZoneFile(){
	//A very raw way of configuring string input to file, need to make this better
	insert_string := `
	,"`+strconv.Itoa(*deviceIDPtr)+"."+*domainNamePtr+`": { "alias": "regular" }
  }
}
`


	f, err := os.OpenFile(DnsZoneFilePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Print("error : "+err.Error())
	}
	defer f.Close()

	newPos, err := f.Seek(-7, 2); if err != nil {
		fmt.Print("error : "+err.Error())
	}

	if _, err := f.WriteAt([]byte(insert_string), newPos); err != nil {
		fmt.Print("error : "+err.Error())
	}

}



func changeDNSRoute(){

	input, err := ioutil.ReadFile(DnsZoneFilePath)
	if err != nil {
		log.Fatalln(err)
	}
fmt.Println(DnsZoneFilePath)

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if *route == "premium"{
			if strings.Contains(line, strconv.Itoa(*deviceIDPtr)+"."+*domainNamePtr ) {
//				fmt.Println("line number is "+ strconv.Itoa(i))
				lines[i] = strings.Replace( lines[i], "regular","premium", -1)
				break
			}
		}
		if *route == "regular"{
//			fmt.Println(lines[i])
			if strings.Contains( line, strconv.Itoa(*deviceIDPtr)+"."+*domainNamePtr ) {
				lines[i] = strings.Replace( lines[i], "premium","regular", -1)
				break
			}
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(DnsZoneFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeDNSEntry(){

	input, err := ioutil.ReadFile(DnsZoneFilePath)
	if err != nil {
		log.Fatalln(err)
	}
//	fmt.Println(DnsZoneFilePath)

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {


			if strings.Contains( line, strconv.Itoa(*deviceIDPtr)+"."+*domainNamePtr ) {
				lines[i] = ""
				break
			}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(DnsZoneFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeDNSEntryForService(){

	input, err := ioutil.ReadFile(DnsZoneFilePath)
	if err != nil {
		log.Fatalln(err)
	}
//	fmt.Println(DnsZoneFilePath)

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {


		if strings.Contains( line, strconv.Itoa(*deviceIDPtr)+"." ) {
			lines[i] = ""
//			fmt.Println(*deviceIDPtr)

		}
	}

		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(DnsZoneFilePath, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}

}



func suspendCDN(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, _ := range lines {
		//		j, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
		//		k, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)


			j, _ := lineNumForString(i,"##rewrite_suspend", TemplatePath)
//		fmt.Println(j)
		if j != 0{
			lines[j+1] = "\tlocation ~* { root /usr/local/nginx/html ; try_files $uri /suspend.html break ;} "
		}



		/*		for looper := j; looper <= k; looper++ {
				lines[looper] = "#" + lines[looper]
				}
				break*/

	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func unsuspendCDN(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, _ := range lines {
		//		j, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
		//		k, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)

//		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##rewrite_suspend", TemplatePath)
		if j != 0 {
			lines[j+1] = ""
		}
//			break
//		}

		/*		for looper := j; looper <= k; looper++ {
				lines[looper] = "#" + lines[looper]
				}
				break*/

	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func terminateCDN(){

	err :=  os.Rename(TemplatePath, TerminatePath+".terminate")

	if err != nil {
//		fmt.Print("error : "+err.Error())
		return
	}

	removeDNSEntryForService()
}



func domainLevelSuspendCDN(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
//		j, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
//		k, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)

		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##rewrite_suspend", TemplatePath)
			lines[j+1] = "\tlocation ~* { root /usr/local/nginx/html ; try_files $uri /suspend.html break ;} "
			break
		}

/*		for looper := j; looper <= k; looper++ {
		lines[looper] = "#" + lines[looper]
		}
		break*/

	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func domainLevelUnsuspendCDN(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {

		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##rewrite_suspend", TemplatePath)
			lines[j+1] = ""
			break
		}

	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func domainLevelTerminateCDN(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, _ := range lines {
		j, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
		k, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)

		for looper := j; looper <= k; looper++ {
			lines[looper] = ""
//			liness := bytes.Replace([]byte(lines[looper]), []byte("substring\n"), []byte(""), 1)
//			lines[looper] = string(liness)

		}
		break

	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	removeDNSEntry()

}



func activateBotnet(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##botnet", TemplatePath)
			lines[j+1] = "\tinclude /usr/local/nginx/conf/aes.conf;\tinclude /usr/local/nginx/conf/bots.conf; "
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func deactivateBotnet(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##botnet", TemplatePath)
			lines[j+1] = ""
			lines[j+2] = ""
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func activateHotlink(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##hotlink", TemplatePath)
			lines[j+1] = "\tinclude /usr/local/nginx/conf/hotlink.conf; "
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func deactivateHotlink(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##hotlink", TemplatePath)
			lines[j+1] = ""
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func enableGzip(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
//	var line_no

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr){
				j, _ := lineNumForString(i,"##gzip", TemplatePath)
//					fmt.Println(j)
					lines[j+1] = "\tgzip on;\tgzip_disable \"msie6\"; \tgzip_vary on;\tgzip_proxied any;\tgzip_comp_level 6;\tgzip_buffers 16 8k;\tgzip_http_version 1.1;\tgzip_types text/html text/css text/plain text/xml text/x-component text/javascript application/x-javascript application/javascript application/json application/manifest+json application/xml application/xhtml+xml application/rss+xml application/atom+xml application/vnd.ms-fontobject application/x-font-ttf application/x-font-opentype application/x-font-truetype image/svg+xml image/x-icon image/vnd.microsoft.icon font/ttf font/eot font/otf font/opentype ;"
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func disableGzip(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##gzip", TemplatePath)
//			fmt.Println(j)
				lines[j+1] = ""

			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func addGzipMimeType(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr){
			j, _ := lineNumForString(i,"##gzip", TemplatePath)
//			fmt.Println(j)
//			lines[j+1] = strings.Replace(lines[j+8], ";", " "+*gzipCustomMimeTypePtr, -1) + " ;"
			str := strings.SplitAfter(lines[j+1], ";")
			str[len(str) - 2] = strings.Replace(str[len(str) - 2],";", "", -1 ) + " "+*gzipCustomMimeTypePtr+" ;"
			lines[j+1] = strings.Join(str, "")
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeGzipMimeType(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr){
			j, _ := lineNumForString(i,"##gzip", TemplatePath)
//			fmt.Println(j)
			lines[j+1] = strings.Replace(lines[j+1],  *gzipCustomMimeTypePtr, "", -1)
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func gzipLocationMimeExclusion(location_path string, mimetype string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = makeGzipExclusion(k, line_no_for_global_gzip, "##gzip", mimetype, lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)


				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = makeGzipExclusion(m, line_no_for_global_gzip, "##gzip", mimetype, lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func makeGzipExclusion(k int, line_no_for_global_gzip int, insert_loc string, mimetype string, lines []string) []string {

		j, _ := lineNumForString(k, insert_loc, TemplatePath)

		//	fmt.Println(k, j)

		if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
		//		fmt.Println("This is a test...")
		lines[j+1] = lines[line_no_for_global_gzip+1]
		lines[j+1] = strings.Replace(lines[j+1],  mimetype, "", -1)
		}else{
		lines[j+1] = strings.Replace(lines[j+1],  mimetype, "", -1)
		}

		//	fmt.Println("The value of j is : "+ strconv.Itoa(j))

		return lines

}



func removeGzipLocationMimeExclusion(location_path string, mimetype string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = removeGzipExclusion(k, line_no_for_global_gzip, "##gzip", mimetype, lines)

				break

			}else//new location
			{

				fmt.Println("error : Location does not exist...")
				return
				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)


				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = removeGzipExclusion(m, line_no_for_global_gzip, "##gzip", mimetype, lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeGzipExclusion(k int, line_no_for_global_gzip int, insert_loc string, mimetype string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

//	fmt.Println(k, j)

//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
		//		fmt.Println("This is a test...")
//		lines[j+1] = lines[line_no_for_global_gzip+1]
		lines[j+1] = strings.Replace(lines[j+1],  ";", mimetype+" ;", -1)
//	}else{
//		lines[j+1] = strings.Replace(lines[j+1],  mimetype, "", -1)
//	}

	//	fmt.Println("The value of j is : "+ strconv.Itoa(j))

	return lines

}



func makeCacheGlobalEntry(){

	input, err := ioutil.ReadFile(CacheIncludeFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	switch *cacheTimeUnitPtr{

	case "milliseconds" :
		*cacheTimeUnitPtr = "ms"
	case "seconds" :
		*cacheTimeUnitPtr = "s"
	case "minutes" :
		*cacheTimeUnitPtr = "m"
	case "hours" :
		*cacheTimeUnitPtr = "h"
	case "days" :
		*cacheTimeUnitPtr = "d"
	case "weeks" :
		*cacheTimeUnitPtr = "w"
	case "months" :
		*cacheTimeUnitPtr = "M"
	case "years" :
		*cacheTimeUnitPtr = "y"
	default :
		fmt.Println("error - Invalid time unit input")
		return
	}

	lines := strings.Split(string(input), "\n")
	var found_domain bool = false
	for i, line := range lines {
		if strings.Contains(line, *sitenamePtr) {

//			fmt.Println("The value of i is : " + strconv.Itoa(i))
			lines[i] = strings.Replace(lines[i], lines[i],
				"proxy_cache_path " + Path_To_Cache_Directory+*sitenamePtr+
					" levels=1:2 keys_zone=" + *sitenamePtr+ ":10m max_size=1g inactive=" + *cacheTimePtr + *cacheTimeUnitPtr+ " use_temp_path=off;", -1)

			found_domain = true

			break
		}

	}

	if found_domain == false{

		lines[len(lines)-1] = "proxy_cache_path " +Path_To_Cache_Directory+*sitenamePtr+
			" levels=1:2 keys_zone="+*sitenamePtr+":10m max_size=1g inactive="+*cacheTimePtr+*cacheTimeUnitPtr+" use_temp_path=off;\n"

	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(CacheIncludeFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func bypassCacheForLocation(location_path string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = makeCacheBypass(k, "##cache_bypass", lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = makeCacheBypass(m, "##cache_bypass", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func makeCacheBypass(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

//	fmt.Println(k, j)

//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
		//		fmt.Println("This is a test...")
		lines[j+1] = lines[j+1] + "proxy_cache_bypass $cookie_nocache $arg_nocache;"

//	}

	//	fmt.Println("The value of j is : "+ strconv.Itoa(j))

	return lines

}



func removeBypassCacheForLocation(location_path string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = removeCacheBypass(k, "##cache_bypass", lines)

				break

			}else//new location
			{

				fmt.Println("error : Location does not exist...")
				return
				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = removeCacheBypass(m, "##cache_bypass", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeCacheBypass(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	//	fmt.Println(k, j)

	//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
	//		fmt.Println("This is a test...")
	lines[j+1] = strings.Replace(lines[j+1] , "proxy_cache_bypass $cookie_nocache $arg_nocache;", "", -1)
//	fmt.Println(strings.Index(lines[j+1], ";"))
	if strings.Index(strings.Trim(strings.Trim(lines[j+1], " "), "\t") , ";") == 0{
		lines[j+1] = ""
	}

	return lines

}



func makeCacheTTLForLocation(location_path string, ttlTime string, ttlTimeUnit string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	switch ttlTimeUnit{

	case "milliseconds" :
		ttlTimeUnit = "ms"
	case "seconds" :
		ttlTimeUnit = "s"
	case "minutes" :
		ttlTimeUnit = "m"
	case "hours" :
		ttlTimeUnit = "h"
	case "days" :
		ttlTimeUnit = "d"
	case "weeks" :
		ttlTimeUnit = "w"
	case "months" :
		ttlTimeUnit = "M"
	case "years" :
		ttlTimeUnit = "y"
	default :
		fmt.Println("error - Invalid time unit input")
		return
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = makeTTL(k, "##cache", ttlTime, ttlTimeUnit, lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = makeTTL(m, "##cache", ttlTime, ttlTimeUnit, lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func makeTTL(k int, insert_loc string, ttlTime string, ttlTimeUnit string,  lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	if !strings.Contains(lines[j+1], "proxy_cache" ){
		lines = makeCacheProxy(k, "##cache", lines)
	}

	if strings.Contains(lines[j+1], "  proxy_cache_valid any " ){

		lines[j+1] = strings.Replace(lines[j+1], substring( lines[j+1], "proxy_cache_valid any ", ";" ), "", -1)


	}else{
		lines[j+1] = lines[j+1] + "  proxy_cache_valid any "+ttlTime+ttlTimeUnit+" ;"
	}

	//	}

	//	fmt.Println("The value of j is : "+ strconv.Itoa(j))

	return lines

}



func removeCacheTTLForLocation(location_path string, ttlTime string, ttlTimeUnit string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	switch ttlTimeUnit{

	case "milliseconds" :
		ttlTimeUnit = "ms"
	case "seconds" :
		ttlTimeUnit = "s"
	case "minutes" :
		ttlTimeUnit = "m"
	case "hours" :
		ttlTimeUnit = "h"
	case "days" :
		ttlTimeUnit = "d"
	case "weeks" :
		ttlTimeUnit = "w"
	case "months" :
		ttlTimeUnit = "M"
	case "years" :
		ttlTimeUnit = "y"
	default :
		fmt.Println("error - Invalid time unit input")
		return
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = removeTTL(k, "##cache", ttlTime, ttlTimeUnit, lines)

				break

			}else//new location
			{

				fmt.Println("error : Location does not exist...")
				return
				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = removeTTL(m, "##cache", ttlTime, ttlTimeUnit, lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeTTL(k int, insert_loc string, ttlTime string, ttlTimeUnit string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

//	fmt.Println(k, j)

	//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
//	fmt.Println("This is a test...")

//	lines[j+1] = strings.Replace(lines[j+1], "proxy_cache_valid any "+ttlTime+ttlTimeUnit+" ;", "", -1)
//	fmt.Println(substring(lines[j+1], "proxy_cache_valid any" , ";"))
	lines[j+1] = strings.Replace(lines[j+1], "proxy_cache_valid any "+ttlTime+ttlTimeUnit+" ;", "", -1)
	//	fmt.Println(strings.Index(lines[j+1], ";"))
	if strings.Index(strings.Trim(strings.Trim(lines[j+1], " "), "\t") , ";") == 0{
		lines[j+1] = ""
	}

	return lines

}



func enableStaleContentForLocation(location_path string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = enableStaleContent(k, "##cache", lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = enableStaleContent(m, "##cache", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func enableStaleContent(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	//	fmt.Println(k, j)

	//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
	//		fmt.Println("This is a test...")
	if !strings.Contains(lines[j+1], "proxy_cache" ){
		lines = makeCacheProxy(k, "##cache", lines)
	}

	lines[j+1] = lines[j+1] + " proxy_cache_use_stale error timeout http_500 http_502 http_503 http_504;"

	//	}

	//	fmt.Println("The value of j is : "+ strconv.Itoa(j))

	return lines

}



func disableStaleContentForLocation(location_path string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = disableStaleContent(k, "##cache", lines)

				break

			}else//new location
			{

				fmt.Println("error : Location does not exist...")
				return
				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = disableStaleContent(m, "##cache", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func disableStaleContent(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	//	fmt.Println(k, j)

	//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
	//		fmt.Println("This is a test...")
	lines[j+1] = strings.Replace(lines[j+1] ,
		" proxy_cache_use_stale error timeout http_500 http_502 http_503 http_504;", "", -1)
	//	fmt.Println(strings.Index(lines[j+1], ";"))
	if strings.Index(strings.Trim(strings.Trim(lines[j+1], " "), "\t") , ";") == 0{
		lines[j+1] = ""
	}

	return lines

}



func enableStringQueryCacheForLocation(location_path string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = enableStringQueryCache(k, "##cache", lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = enableStringQueryCache(m, "##cache", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func enableStringQueryCache(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)


	if !strings.Contains(lines[j+1], "proxy_cache" ){
		lines = makeCacheProxy(k, "##cache", lines)
	}

	lines[j+1] = lines[j+1] + "proxy_cache_key $host$uri$is_args$args;"

	//	}

	//	fmt.Println("The value of j is : "+ strconv.Itoa(j))

	return lines

}



func disableStringQueryCacheForLocation(location_path string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = disableStringQueryCache(k, "##cache", lines)

				break

			}else//new location
			{

				fmt.Println("error : Location does not exist...")
				return
				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = disableStringQueryCache(m, "##cache", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func disableStringQueryCache(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	//	fmt.Println(k, j)

	//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
	//		fmt.Println("This is a test...")
	lines[j+1] = strings.Replace(lines[j+1] , "proxy_cache_key $host$uri$is_args$args;", "proxy_cache_key $host$uri;", -1)
	//	fmt.Println(strings.Index(lines[j+1], ";"))
	if strings.Index(strings.Trim(strings.Trim(lines[j+1], " "), "\t") , ";") == 0{
		lines[j+1] = ""
	}

	return lines

}



func removeCacheForLocation(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {

		if strings.Contains(line, "##"+*domainNamePtr) {

			k, _ := lineNumForString(i, *remove_cache_location, TemplatePath)
//			fmt.Println(k)
			if k!= 0 {
				j, _ := lineNumForString(k, "##cache", TemplatePath)
				lines[j+1] = ""
				break
			}

		}

	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func enableAllCacheOptions(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	//	var line_no

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr){
			j, _ := lineNumForString(i,"##cache", TemplatePath)
			//					fmt.Println(j)

			if strings.Contains(lines[j+1], "##") {
				strings.Replace(lines[j+1], "##","", 1)
			}
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func disableAllCacheOptions(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	//	var line_no

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr){
			j, _ := lineNumForString(i,"##cache", TemplatePath)
			//					fmt.Println(j)

				lines[j+1] = "##"+lines[j+1]

			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func purgeByValue(purge_method string, purge_val string) {
	searchDir := "C:\\Users\\Getit\\Desktop\\cache2"+*sitenamePtr

	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)

		if !f.IsDir() {
			//			fmt.Println(path+f.Name())
			readFileByCachePath(path, "*.png")
		}

		return nil
	})

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
			if strings.Contains(scanner.Text(), search_str){
				fmt.Println(scanner.Text())
				//				os.Remove(path)
				break
			}
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}



func purgeCache(){

	switch *purgeMethodPtr{

	case "entirecache" :
		out, _ := exec.Command("/root/purge_request.sh", *purgeMethodPtr, *sitenamePtr).Output()
		fmt.Println(out)

	case "filename" :
		out, _ := exec.Command("/root/purge_request.sh", *purgeMethodPtr, *purgeFileNamePtr, *sitenamePtr).Output()
		fmt.Println(out)

	case "extension" :
		out, _ := exec.Command("/root/purge_request.sh", *purgeMethodPtr, *purgeFileExtPtr, *sitenamePtr).Output()
		fmt.Println(out)

	case "path" :
		out, _ := exec.Command("/root/purge_request.sh", *purgeMethodPtr, *purgePathPtr, *sitenamePtr).Output()
		fmt.Println(out)

	default :
		fmt.Println("error - Invalid purgemethod input")
		return
	}

}



func enableAutoSSL(){

//	exec.Command("certbot certonly --webroot --agree-tos --no-eff-email --email kamalkishor.pande@psychz.net -w /var/www/letsencrypt -d www."+*domainNamePtr+" -d "+*domainNamePtr).Output()
	out, _ := exec.Command("certbot", "certonly", "--webroot", "--agree-tos", "--no-eff-email", "--email", "kamalkishor.pande@psychz.net", "-w", "/var/www/letsencrypt", "-d", "www."+*domainNamePtr, "-d"+*domainNamePtr).Output()
	fmt.Printf("%s", out)

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))
			j, _ := lineNumForString(i,"##ssllisten", TemplatePath)
			lines[j+1] = "\t\tlisten 172.106.22.7:443 ;"
			//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")

			k, _ := lineNumForString(i,"##sslconf", TemplatePath)
			lines[k+1] = "\t\tssl_certificate /etc/letsencrypt/live/www."+*domainNamePtr+"/fullchain.pem; " +
				"ssl_certificate_key /etc/letsencrypt/live/www."+*domainNamePtr+"/privkey.pem; " +
				"ssl_trusted_certificate /etc/letsencrypt/live/www."+*domainNamePtr+"/fullchain.pem; " +
				"include /usr/local/nginx/conf/snippets/ssl.conf ;"+
			" if ( $scheme = \"http\" ) { return 301 https://www.$server_name$request_uri; } "
			//add if statement for scheme

			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func disableAutoSSL(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))
			j, _ := lineNumForString(i,"##ssllisten", TemplatePath)
			lines[j+1] = ""
			//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")

			k, _ := lineNumForString(i,"##sslconf", TemplatePath)
			lines[k+1] = ""
			//+"\t\tssl_certificate /usr/local/nginx/ssl/"+*domainNamePtr+"/fullchain.pem; " +
			//				"ssl_certificate_key /usr/local/nginx/ssl/"+*domainNamePtr+"/privkey.pem; " +
			//				"ssl_trusted_certificate /etc/letsencrypt/live/www.mydomain.com/fullchain.pem; " +
			//				"include /etc/nginx/ssl/snippets/ssl.conf ;"

			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func enableThirdPartySSL(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))
			j, _ := lineNumForString(i,"##ssllisten", TemplatePath)
			lines[j+1] = "\t\tlisten 172.106.22.7:443 ;" + " if ( $scheme = \"http\" ) { return 301 https://$server_name$request_uri; } "
			//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")

			k, _ := lineNumForString(i,"##sslconf", TemplatePath)
			lines[k+1] = "\t\tssl_certificate /usr/local/nginx/conf/SSL/"+*sitenamePtr+"/"+*sitenamePtr+".crt; " +
							"ssl_certificate_key /usr/local/nginx/conf/SSL/"+*sitenamePtr+"/"+*sitenamePtr+".key; " //+
//							"ssl_trusted_certificate /usr/local/nginx/conf/SSL/"+*sitenamePtr+"/"+*sitenamePtr+".crt; "

			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func disableThirdPartySSL(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))
			j, _ := lineNumForString(i,"##ssllisten", TemplatePath)
			lines[j+1] = ""
			//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")

			k, _ := lineNumForString(i,"##sslconf", TemplatePath)
			lines[k+1] = ""

			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func enableLayer7(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))

			j, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
			k, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)

			for looper := j; looper <= k; looper++ {
				l, _ := lineNumForString(looper,"##layer_7", TemplatePath)
				if (!strings.Contains(lines[l+1], "#")) && (l < k) {
					//fmt.Println(l+1)
					lines[l+1] = "\t\ttestcookie on ; "
					i = l

				}
				if (strings.Contains(lines[l+1], "#")) && (l < k) {
				//	fmt.Println(lines[l+1])
					lines[l+1] = lines[l+1][1:]
					i = l

				}
			}

			break

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func disableLayer7(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))

			j, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
			k, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)

			for looper := j; looper <= k; looper++ {
				l, _ := lineNumForString(looper,"##layer_7", TemplatePath)
				if l < k {
					//fmt.Println(l+1)
					if !strings.Contains(lines[l+1], "#"){
						lines[l+1] = "#" + lines[l+1]	//this section of the code needs tuning - the entire nested loop
					}

				}

			}
			break

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func layer7ExemptionPathEnable(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, *layer7LocationPtr, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = addLayer7ExemptionPath(k, "##layer_7", lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, *layer7LocationPtr)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = addLayer7ExemptionPath(m, "##layer_7", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func addLayer7ExemptionPath(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	lines[j+1] = "\ttestcookie off; "

	return lines

}



func layer7ExemptionPathDisable(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {

			//			line_no_for_global_gzip, _ := lineNumForString(i,"##gzip", TemplatePath)

			k, _ := lineNumForString(i, *layer7LocationPtr, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = removeLayer7ExemptionPath(k, "##layer_7", lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, *layer7LocationPtr)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = removeLayer7ExemptionPath(m, "##layer_7", lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeLayer7ExemptionPath(k int, insert_loc string, lines []string) []string {

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	lines[j+1] = "\ttestcookie on; "

	return lines

}



func layer7WhiteListAdd(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(strings.Trim(strings.Trim(string(input), "\t"), ""), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))
			j, _ := lineNumForString(i,"##layer_7", TemplatePath)

			//o, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
			p, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)

			if j < p {
				if !strings.Contains(lines[j+1], "testcookie_whitelist") {

					lines[j+1] = lines[j+1] + " testcookie_whitelist { " + *ip_whitelistPtr + " ; }"

					break
				} else
				{

					lines[j+1] = strings.Replace(lines[j+1], "}", *ip_whitelistPtr+" ; }", -1)
					break
				}
			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func layer7WhiteListRemove(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(strings.Trim(strings.Trim(string(input), "\t"), ""), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of i is : "+ strconv.Itoa(i))
			j, _ := lineNumForString(i,"##layer_7", TemplatePath)
			//o, _ := lineNumForString(i,"##"+*domainNamePtr, TemplatePath)
			p, _ := lineNumForString(i,"##"+*domainNamePtr+"_ends", TemplatePath)

			if j < p {
				if strings.Contains(lines[j+1], *ip_whitelistPtr) {

					lines[j+1] = strings.Replace(lines[j+1], *ip_whitelistPtr+" ;", "", -1)
					if !strings.Contains(lines[j+1], "."){
						lines[j+1] = "\t\ttestcookie on ; "
					}

					break
				}
			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



//func createGzipLocation(){
//
//	input, err := ioutil.ReadFile(TemplatePath)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	lines := strings.Split(string(input), "\n")
//
//	for i, line := range lines {
//		if strings.Contains(line, "##"+*domainNamePtr) {
//			k, _ := lineNumForString(i, *gzipLocationPtr, TemplatePath)
//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
//			j, _ := lineNumForString(i,"##location / ends", TemplatePath)
//			lines[j+1] = "\n\t\t##location " + *gzipLocationPtr + " starts" + "\n\t\tlocation "+*gzipLocationPtr+"\n\t\t{\n\t\t\n\t\t" +
//				"\tgzip on;\n\tgzip_disable \"msie6\";\n\tgzip_vary on;\n\tgzip_proxied any;\n\tgzip_comp_level 6;\n\tgzip_buffers 16 8k;\n\tgzip_http_version 1.1;\n\tgzip_types text/html text/css text/plain text/xml text/x-component text/javascript application/x-javascript application/javascript application/json application/manifest+json application/xml application/xhtml+xml application/rss+xml application/atom+xml application/vnd.ms-fontobject application/x-font-ttf application/x-font-opentype application/x-font-truetype image/svg+xml image/x-icon image/vnd.microsoft.icon font/ttf font/eot font/otf font/opentype ; \n" +
//				"\n\n\t\t}\n\t\t##location "+*gzipLocationPtr+" ends"
//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")
//			break
//		}
//	}
//	output := strings.Join(lines, "\n")
//	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//}



func configParkedDomain(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		//This ensures that the parked domain is appended to the said domain only
		if strings.Contains(line, "##"+*domainNamePtr){
			j, _ := lineNumForString(i,"##server_name", TemplatePath)
//			fmt.Println(j)
			lines[j+1] = strings.Replace(lines[j+1],  ";", "  "+*parkeddomainPtr + " ;", -1)
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeParkedDomain(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
//	var final_line string

	for i, line := range lines {
		//This ensures that the parked domain is appended to the said domain only
		if strings.Contains(line, "##"+*domainNamePtr){
			j, _ := lineNumForString(i,"##server_name", TemplatePath)
//			fmt.Println(j)

			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), *parkeddomainPtr)//final_line//strings.Replace(lines[j+1], *parkeddomainPtr, "", -1)
//			lines[i+1] = strings.Replace(lines[i+1], "", "", -1) + " ;"
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func configBandwidth()  {

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##bandwidth", TemplatePath)
			lines[j+1] = "\t\tlimit_rate " + strconv.Itoa(*bandwidthPtr) + string(*bandwidthUnitPtr)+" ;"
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeBandwidthLimit()  {

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			j, _ := lineNumForString(i,"##bandwidth", TemplatePath)
			lines[j+1] = ""
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func createPageSpeed(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
//			k, _ := lineNumForString(i,location_path, file_path)
//			fmt.Println("The value of i is : "+ strconv.Itoa(i))
			j, _ := lineNumForString(i,"##pagespeed", TemplatePath)
			lines[j+1] = "\tinclude /usr/local/nginx/conf/pagespeed.conf ;"
//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removePageSpeed(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			j, _ := lineNumForString(i,"##pagespeed", TemplatePath)
			lines[j+1] = ""
			if strings.Contains(lines[j+2], "pagespeed Disallow"){
				lines[j+2] = ""
			}
			//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func blockPageSpeedByLocation(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			j, _ := lineNumForString(i, "##pagespeed", TemplatePath)
			if lines[j+1] == "\tinclude /usr/local/nginx/conf/pagespeed.conf ;" {
				lines[j+2] = lines[j+2] + "\tpagespeed Disallow \"*/" + *pagespeedLocationPtr + "/*\" ; \t"
			//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")
				break
			}
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func unblockPageSpeedByLocation(){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			//			k, _ := lineNumForString(i,location_path, file_path)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			j, _ := lineNumForString(i, "##pagespeed", TemplatePath)
			if lines[j+1] == "\tinclude /usr/local/nginx/conf/pagespeed.conf ;" {
				lines[j+2] = strings.Replace(lines[j+2], "pagespeed Disallow \"*/" + *pagespeedLocationPtr + "/*\" ; ", "", -1)
				//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")
				break
			}
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func addAccessControlAtLocation(location_path string, iptoallow string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			k, _ := lineNumForString(i, location_path, TemplatePath)
//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = make_access_control(k, "##access_control", iptoallow, lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

//				forceWrite(lines)
//				l, _ := lineNumForString(m, location_path, TemplatePath)
//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = make_access_control(m, "##access_control", iptoallow, lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func make_access_control(k int, insert_loc string, iptoallow string, lines []string) []string{

	j, _ := lineNumForString(k,insert_loc, TemplatePath)
//	fmt.Println("The value of j is : "+ strconv.Itoa(j))
	if(*accessControlWhitelistIPPtr!="") {
		lines[j+1] = strings.Replace(lines[j+1], "deny all;", "", -1)
		lines[j+1] = strings.Replace(lines[j+1], "allow all;", "", -1)
		lines[j+1] = lines[j+1] + "\tallow " + iptoallow + "; deny all;"
	}
	if(*accessControlBlacklistIPPtr!=""){
		lines[j+1] = strings.Replace(lines[j+1], "allow all;", "", -1)
		lines[j+1] = strings.Replace(lines[j+1], "deny all;", "", -1)
		if strings.Contains(lines[j+1], "allow"){
			lines[j+1] = lines[j+1] + "\tdeny " + iptoallow + "; deny all;"
		}else{
			lines[j+1] = lines[j+1] + "\tdeny " + iptoallow + "; allow all;"
		}
	}

	return lines

}



func removeWhiteListAtLocation(location string, iptoremove string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			k, _ := lineNumForString(i, location, TemplatePath)
//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			if k != 0{

				j, _ := lineNumForString(i,"##access_control", TemplatePath)
				if(*accessControlWhitelistIPPtr!="") {
					lines[j+1] = strings.Replace(lines[j+1], "deny all;", "", -1)
					lines[j+1] = strings.Replace(lines[j+1], "allow all;", "", -1)
					if strings.Contains(lines[j+1], "allow " + iptoremove+ ";"){
						lines[j+1] = strings.Replace(lines[j+1], "allow " + iptoremove+ ";", "", -1)
					}
					if strings.Contains(lines[j+1], "allow") && strings.Contains(lines[j+1], "deny"){
						lines[j+1] = lines[j+1] + " deny all;"
					}
					if !strings.Contains(lines[j+1], "allow"){
						lines[j+1] = lines[j+1] + " allow all;"
					}
					if !strings.Contains(lines[j+1], "deny"){
						lines[j+1] = lines[j+1] + " deny all;"
					}
					if !strings.Contains(lines[j+1], "."){
						lines[j+1] = ""
					}

//					lines[j+1] = lines[j+1] + " allow   " + iptoremove + "; deny all;"
				}

				//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")
				break

			}else
			{

				fmt.Println("error : Location does not exist, please check...")
				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func removeBlackListAtLocation(location string, iptoremove string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			k, _ := lineNumForString(i, location, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			if k != 0{

				j, _ := lineNumForString(i,"##access_control", TemplatePath)
				if(*accessControlBlacklistIPPtr!="") {
					lines[j+1] = strings.Replace(lines[j+1], "deny all;", "", -1)
					lines[j+1] = strings.Replace(lines[j+1], "allow all;", "", -1)
					if strings.Contains(lines[j+1], "deny " + iptoremove+ ";"){
						lines[j+1] = strings.Replace(lines[j+1], "deny " + iptoremove+ ";", "", -1)
					}
					if strings.Contains(lines[j+1], "allow") && strings.Contains(lines[j+1], "deny"){
						lines[j+1] = lines[j+1] + " deny all;"
					}
					if !strings.Contains(lines[j+1], "allow"){
						lines[j+1] = lines[j+1] + " allow all;"
					}
					if !strings.Contains(lines[j+1], "deny"){
						lines[j+1] = lines[j+1] + " deny all;"
					}
					if !strings.Contains(lines[j+1], "."){
						lines[j+1] = ""
					}

//					lines[j+1] = lines[j+1] + " allow   " + iptoremove + "; deny all;"
				}

				//			lines[j+1] = stringReplaceStrictMatch(string(lines[j+1]), "font/opentype")
				break

			}else
			{

				fmt.Println("error : Location does not exist, please check...")
				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}


//For cache ON/OFF feature
func addCacheControlAtLocation(location_path string, iptoallow string){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "##"+*domainNamePtr) {
			k, _ := lineNumForString(i, location_path, TemplatePath)
			//			fmt.Println("The value of k is : "+ strconv.Itoa(k))
			//existing location
			if k != 0{

				lines = make_cache_control(k, "##cache", iptoallow, lines)

				break

			}else//new location
			{

				m, _ := lineNumForString(i,"##custom_location", TemplatePath)

				//				fmt.Println("value of m is : "+ strconv.Itoa(m))
				lines = create_new_blank_location(m, location_path)

				//				l, _ := lineNumForString(m, location_path, TemplatePath)
				//				fmt.Println("value of l is : "+ strconv.Itoa(l))
				lines = make_cache_control(m, "##cache", iptoallow, lines)

				break

			}

		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}


//For cache ON/OFF feature
func make_cache_control(k int, insert_loc string, iptoallow string, lines []string) []string{

	j, _ := lineNumForString(k,"##cache", TemplatePath)
	//	fmt.Println("The value of j is : "+ strconv.Itoa(j))
	if(*accessControlWhitelistIPPtr!="") {
		lines[j+1] = strings.Replace(lines[j+1], "deny all;", "", -1)
		lines[j+1] = strings.Replace(lines[j+1], "allow all;", "", -1)
		lines[j+1] = lines[j+1] + "\tallow " + iptoallow + "; deny all;"
	}
	if(*accessControlBlacklistIPPtr!=""){
		lines[j+1] = strings.Replace(lines[j+1], "allow all;", "", -1)
		lines[j+1] = strings.Replace(lines[j+1], "deny all;", "", -1)
		if strings.Contains(lines[j+1], "allow"){
			lines[j+1] = lines[j+1] + "\tdeny " + iptoallow + "; deny all;"
		}else{
			lines[j+1] = lines[j+1] + "\tdeny " + iptoallow + "; allow all;"
		}
	}

	return lines

}



func createFile() {
	// detect if file exists
	var _, err = os.Stat(TemplatePath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(TemplatePath)
		checkError(err) //okay to call os.exit()
		defer file.Close()
	}
}



func writeFile() {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(TemplatePath, os.O_RDWR, 0644)
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



func readFile()(string) {
	// re-open file
	file, err := os.OpenFile(TemplatePath, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// read file
	var text = make([]byte, 1024)
	n, err := file.Read(text)
	if n > 0 {
		return string(text)
	}
	//if there is an error while reading
	//just print however much was read if any
	//at return file will be closed
	return ""
}



func deleteFile() {
	// delete file
	var err = os.Remove(TemplatePath)
	checkError(err)
}



func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}



func checkIfFileExists()(exists bool){
	if _, err := os.Stat(TemplatePath); !os.IsNotExist(err) {
		return true
	}
	return false
}



func lineNumForString(found_at int, search_string string, path string)(int, error){

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for j, line := range lines[found_at:]{

//		if j==0 {
//			fmt.Println("The line and search string are : .... "+line+"  ----  "+search_string)
//		}

		var words []string
		if strings.Contains(line, search_string) {
			words = strings.Split(line, " ")
		}

		for _, word := range words{

//			fmt.Println(strconv.Itoa(i) + " --word number " + word+" --line number "+strconv.Itoa(j+found_at) +" search string is : "+strings.Trim(word, "\t"))
			if strings.Trim(strings.Trim(word, "\t"), " ") == search_string {
//				if word ==
//				fmt.Println(strconv.Itoa(i) + " --word number " + word+" --line number "+strconv.Itoa(j+found_at) +" search string is : "+search_string)

				return j+found_at, err
			}
		}

	}

	return 0, err

}



func stringReplaceStrictMatch(original_string string, to_remove_value string) string{


	words := strings.Split(original_string, " ")
	var replaced_string string

	//			fmt.Println("This is a test : " + string(lines[i]) )
	for k, word := range words {

//		fmt.Println(strconv.Itoa(k) + " " + word)

		if to_remove_value == word{

			words[k] = ""
		}
		replaced_string = strings.Join(words, " ")

	}

	return replaced_string
}



func findStringExactMatch(original_string string, to_search_value string) bool{


	words := strings.Split(original_string, " ")
	//var replaced_string string

	for _, word := range words {

		fmt.Println(word)
		if to_search_value == word{

			return true
		}else{
			return false
		}
	}
	return false
	//return replaced_string
}



func makeCacheProxy(k int, insert_loc string, lines []string) []string{

	j, _ := lineNumForString(k, insert_loc, TemplatePath)

	//	fmt.Println(k, j)

	//	if strings.Trim(strings.Trim(lines[j+1], "\t"), "") == "" {
	//		fmt.Println("This is a test...")
	if strings.Contains(lines[j+1], " proxy_cache "+*sitenamePtr+";"){
		return lines
	}else{
		lines[j+1] = lines[j+1] + " proxy_cache "+*sitenamePtr+";"
	}

	//	}

	//	fmt.Println("The value of j is : "+ strconv.Itoa(j))

	return lines

}



func create_new_blank_location(lineno_to_insert_location int, location_path string) []string{

	input, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	httpscheme := "http";

	if *sslOnPtr == "false"{
		httpscheme = "http"
	}

	if *sslOnPtr == "true"{
		httpscheme = "https"
	}

//		k, _ := lineNumForString(i, location_path, TemplatePath)
			lines[lineno_to_insert_location+1] = "\n\n\tlocation " + location_path +
				" {\n" +
				"\t\t##access_control\n\t\t" +
				"\n" +
				"\t\t##layer_7\n\t\t" +
				"\n" +
				"\t\t##gzip\n\t\t" +
				"\n" +
				"\t\t##cache_bypass\n\t\t" +
				"\n" +
				"\t\t##cache\n\t\t" +
				"proxy_redirect     off;"+
				"proxy_set_header   Host			$host;"+
				"proxy_set_header   DeviceID		www."+*domainNamePtr+";"+
				"proxy_set_header   X-Real-IP		$remote_addr;"+
				"proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;"+
				"\n" +
				"\t\tproxy_pass " + httpscheme + "://" +
				*ipAddressPtr + ":" + strconv.Itoa(*domainPortPtr) + ";" +
				"\n\t}" +
				"\n\t\t##custom_location\n"


	forceWrite(lines)

	input2, err := ioutil.ReadFile(TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines2 := strings.Split(string(input2), "\n")

	return lines2

}



func forceWrite(lines []string){

	output := strings.Join(lines, "\n")
	err := ioutil.WriteFile(TemplatePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}



func substring(input string, m string, n string) string{


	var start_index int
	var end_index int
	start_index = strings.Index(input, m)
	end_index = strings.Index(input[start_index:], n)

	return input[start_index:end_index+end_index+1]

}