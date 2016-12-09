package router

import (
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/handler"
)

type Route struct {
	path        string
	handlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{"/skills/", handler.MakeHandler(handler.SkillsHandler, controller.SkillsController{})},
}

func StartRouter() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handlerFunc)
	}
	return
}
