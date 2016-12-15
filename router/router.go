package router

import (
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/handler"
)

/*
Route contains an HTTP URI endpoint (e.g. "/skills" or "/skills/") in the 'path' var, and the
handler function with which to handle HTTP requests to that endpoint in the 'handlerFunc' var.
*/
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
	{
		"/skills/",
		handler.MakeHandler(
			handler.Handler,
			&controller.SkillsController{
				BaseController: &controller.BaseController{},
			})},
	{
		"/skills",
		handler.MakeHandler(
			handler.Handler,
			&controller.SkillsController{
				BaseController: &controller.BaseController{},
			})},
}

/*
StartRouter() instantiates a new http.ServeMux and registers with it each endpoint that is currently being handled
by the SkillDirectory REST API with an appropriate handler function for that endpoint. This http.ServeMux is returned.
*/
func StartRouter() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handlerFunc)
	}
	return
}
