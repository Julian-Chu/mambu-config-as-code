package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Julian-Chu/MambuConfigurationAPI"
)

// `go run . -url=<mambuBaseURL> -apikey=<APIKEY>`
func main() {
	var mambuBaseURL string
	var apikey string
	//var targetArg string
	flag.StringVar(&mambuBaseURL, "url", "", "mambu base url")
	flag.StringVar(&apikey, "apikey", "", "mambu apikey")
	//flag.StringVar(&targetArg, "target", "", "client or account")
	flag.Parse()

	if mambuBaseURL == "" {
		log.Fatalln("missing arg: url")
	}
	if apikey == "" {
		log.Fatalln("missing arg: apikey")
	}
	c := MambuConfigurationAPI.NewClient(mambuBaseURL, apikey)
	fmt.Println(c.GetCustomFields())
}
