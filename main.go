package main

import (
	"net/http"
	"skilldirectory/router"
)

func main() {
	router := router.StartRouter()
	http.ListenAndServe(":8080", router)
}
