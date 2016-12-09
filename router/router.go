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

// To add a new route, add a new to the []Route Slice:
// Route{
// 	"/ENDPOINT/",
// 	handler.MakeHandler(
// 		handler.Handler,
// 		&controller.NEW_CONTROLLER{
// 			BaseController: &controller.BaseController{},
// 		})},
// And add a controller to the controller package

var routes = []Route{
	Route{
		"/skills/",
		handler.MakeHandler(
			handler.Handler,
			&controller.SkillsController{
				BaseController: &controller.BaseController{},
			})},
}

func StartRouter() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handlerFunc)
	}
	return
}
