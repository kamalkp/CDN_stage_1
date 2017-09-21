package main

import (
	"flag"
	"os/exec"
	"fmt"
)


func main(){


	var domainNamePtr = flag.String("domainname", "example.com", "The domain name to be registered without the www")

	flag.Parse()

	s := "certbot certonly --webroot --agree-tos --no-eff-email --email kamalkishor.pande@psychz.net -w /var/www/letsencrypt -d www."+*domainNamePtr+" -d "+*domainNamePtr

	out, _ := exec.Command(s).Output()
	fmt.Printf("%s", out)


}
