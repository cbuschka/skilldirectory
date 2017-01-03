package main

import (
	"flag"
	"fmt"
	"net/http"
	"skilldirectory/data"
	"skilldirectory/router"
)

var session data.DataAccess

var debug bool

func init() {
	flag.Bool("debug", false, "Change log level")
}

func main() {
	flag.Parse()
	debugFlag := flag.Lookup("debug")
	fmt.Println(debugFlag)
	router := router.StartRouter()
	http.ListenAndServe(":8080", router)
}
