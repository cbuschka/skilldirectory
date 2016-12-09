package router

import (
	"net/http"
	"skilldirectory/handler"
)

type Route struct {
	name        string
	path        string
	handlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{"skills", "/skills", handler.MakeHandler(handler.SkillsHandler)},
	Route{"skills", "/skills/", handler.MakeHandler(handler.SkillsHandler)},
}

func StartRouter() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handlerFunc)
	}
	return
}
