package router

import (
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
	"skilldirectory/handler"
	"skilldirectory/helper"
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

var (
	url      = helper.GetProperty("CASSANDRA_URL")
	port     = helper.GetProperty("CASSANDRA_PORT")
	keyspace = helper.GetProperty("CASSANDRA_KEYSPACE")
	session  = data.NewCassandraConnector(url, port, keyspace)
)

var (
	skillsController = controller.SkillsController{
		BaseController: &controller.BaseController{},
	}
	skillsHandlerFunc = handler.MakeHandler(handler.Handler, &skillsController, session)

	teamMembersController = controller.TeamMembersController{
		BaseController: &controller.BaseController{},
	}
	teamMembersHandlerFunc = handler.MakeHandler(handler.Handler, &teamMembersController, session)

	tmSkillsController = controller.TMSkillsController{
		BaseController: &controller.BaseController{},
	}
	tmSkillsHandlerFunc = handler.MakeHandler(handler.Handler, &tmSkillsController, session)

	routes = []Route{
		{"/skills/", skillsHandlerFunc},
		{"/skills", skillsHandlerFunc},
		{"/teammembers/", teamMembersHandlerFunc},
		{"/teammembers", teamMembersHandlerFunc},
		{"/tmskills/", tmSkillsHandlerFunc},
		{"/tmskills", tmSkillsHandlerFunc},
	}
)

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
