package main

import (
	"flag"
	"net/http"
	"skilldirectory/router"
)

var debug bool

func init() {
	flag.Bool("debug", false, "Change log level")
}

func main() {
	flag.Parse()
	router := router.StartRouter()
	http.ListenAndServe(":8080", router)
}
