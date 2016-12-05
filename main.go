package main

import (
	"net/http"
	"skilldir/router"
)

func main() {
	router := router.StartRouter()
	http.ListenAndServe(":8080", router)
}
