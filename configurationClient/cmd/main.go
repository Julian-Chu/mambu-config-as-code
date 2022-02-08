package main

import (
	"flag"
	"fmt"
	"log"

	client "github.com/Julian-Chu/MambuConfigurationAPI/configurationClient/rest"
)

// `go run . -url=<mambuBaseURL> -apikey=<APIKEY>`
func main() {
	var mambuBaseURL string
	var apikey string
	//var targetArg string
	flag.StringVar(&mambuBaseURL, "url", "", "provider base url")
	flag.StringVar(&apikey, "apikey", "", "provider apikey")
	//flag.StringVar(&targetArg, "target", "", "client or account")
	flag.Parse()

	if mambuBaseURL == "" {
		log.Fatalln("missing arg: url")
	}
	if apikey == "" {
		log.Fatalln("missing arg: apikey")
	}
	c := client.NewClient(mambuBaseURL, apikey)
	fmt.Println(c.GetCustomFields())
}
