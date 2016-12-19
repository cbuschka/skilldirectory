package main

import (
	"net/http"
	"skilldirectory/data"
	"skilldirectory/router"
)

var session data.DataAccess

func main() {
	router := router.StartRouter()
	http.ListenAndServe(":8080", router)
}
