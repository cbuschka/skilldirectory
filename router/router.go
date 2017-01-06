package router

import (
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
	"skilldirectory/handler"
	util "skilldirectory/util"
)

/*
Route contains an HTTP URI endpoint (e.g. "/skills" or "/skills/") in the 'path'
var, and the handler function with which to handle HTTP requests to that
endpoint in the 'handlerFunc' var.
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
	url      string
	port     string
	keyspace string
	username string
	password string
	session  *data.CassandraConnector
	routes   []Route
)

// InitializeSession sets global variables at start up
func initializeCassandra() {
	url = util.GetProperty("CASSANDRA_URL")
	port = util.GetProperty("CASSANDRA_PORT")
	keyspace = util.GetProperty("CASSANDRA_KEYSPACE")
	username = util.GetProperty("CASSANDRA_USERNAME")
	password = util.GetProperty("CASSANDRA_PASSWORD")
	session = data.NewCassandraConnector(url, port, keyspace, username, password)
}

func loadRoutes() {

	skillsController := controller.SkillsController{
		BaseController: &controller.BaseController{},
	}
	skillsHandlerFunc := handler.MakeHandler(handler.Handler, &skillsController, session)

	teamMembersController := controller.TeamMembersController{
		BaseController: &controller.BaseController{},
	}
	teamMembersHandlerFunc := handler.MakeHandler(handler.Handler, &teamMembersController, session)

	tmSkillsController := controller.TMSkillsController{
		BaseController: &controller.BaseController{},
	}
	tmSkillsHandlerFunc := handler.MakeHandler(handler.Handler, &tmSkillsController, session)

	linksController := controller.LinksController{
		BaseController: &controller.BaseController{},
	}
	linksHandlerFunc := handler.MakeHandler(handler.Handler, &linksController, session)

	skillReviewsController := controller.SkillReviewsController{
		BaseController: &controller.BaseController{},
	}
	skillReviewsHandlerFunc := handler.MakeHandler(handler.Handler, &skillReviewsController, session)

	routes = []Route{
		{"/skills/", skillsHandlerFunc},
		{"/skills", skillsHandlerFunc},
		{"/teammembers/", teamMembersHandlerFunc},
		{"/teammembers", teamMembersHandlerFunc},
		{"/tmskills/", tmSkillsHandlerFunc},
		{"/tmskills", tmSkillsHandlerFunc},
		{"/links/", linksHandlerFunc},
		{"/links", linksHandlerFunc},
		{"/skillreviews", skillReviewsHandlerFunc},
		{"/skillreviews/", skillReviewsHandlerFunc},
	}
}

/*
StartRouter() instantiates a new http.ServeMux and registers with it each
endpoint that is currently being handled by the SkillDirectory REST API with an
appropriate handler function for that endpoint. This http.ServeMux is returned.
*/
func StartRouter() (mux *http.ServeMux) {
	initializeCassandra()
	loadRoutes()
	mux = http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handlerFunc)
	}
	return
}
