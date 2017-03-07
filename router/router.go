package router

import (
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
	"skilldirectory/handler"
	util "skilldirectory/util"

	"os/user"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
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
	url        string
	port       string
	keyspace   string
	username   string
	password   string
	db         *gorm.DB
	fileSystem data.FileSystem
	routes     []Route
)

// initCassandra sets global variables at start up
func initCassandra() {
	url = util.GetProperty("CASSANDRA_URL")
	port = util.GetProperty("CASSANDRA_PORT")
	keyspace = util.GetProperty("CASSANDRA_KEYSPACE")
	username = util.GetProperty("CASSANDRA_USERNAME")
	password = util.GetProperty("CASSANDRA_PASSWORD")
	url = util.GetProperty("POSTGRES_URL")
	port = util.GetProperty("POSTGRES_PORT")
	keyspace = util.GetProperty("POSTGRES_KEYSPACE")
	username = util.GetProperty("POSTGRES_USERNAME")
	password = util.GetProperty("POSTGRES_PASSWORD")
	db = data.NewPostgresConnector(url, port, keyspace, username, password).DB()
}

// initFileSystem sets global variables at start up
func initFileSystem() {
	fs := util.GetProperty("FILE_SYSTEM")
	switch fs {
	case "S3": // Use AWS S3 as file system
		var err error
		fileSystem, err = data.NewS3Session()
		if err != nil {
			panic("Failed to connect to AWS S3!")
		}
		log.Info("Using AWS S3 as file system.")
	default: // Use local disk as file system by default
		fileSystem = data.NewLocalFileSystem()
		log.Info("Using local disk as file system.")
		user, _ := user.Current()
		log.Infof("Hosting static file server for '%s/skilldirectory' on localhost:2121.",
			user.HomeDir)
		go func() {
			err := http.ListenAndServe(":2121", http.FileServer(http.Dir(
				user.HomeDir+"/skilldirectory")))
			if err != nil {
				log.Errorf("Error produced while running static file server: %s", err)
			}
		}()
	}
}

func loadRoutes() {

	skillsController := controller.SkillsController{
		BaseController: &controller.BaseController{},
	}
	skillsHandlerFunc := handler.MakeHandler(handler.Handler, &skillsController, fileSystem, db)

	teamMembersController := controller.TeamMembersController{
		BaseController: &controller.BaseController{},
	}
	teamMembersHandlerFunc := handler.MakeHandler(handler.Handler, &teamMembersController, fileSystem, db)

	tmSkillsController := controller.TMSkillsController{
		BaseController: &controller.BaseController{},
	}
	tmSkillsHandlerFunc := handler.MakeHandler(handler.Handler, &tmSkillsController, fileSystem, db)

	linksController := controller.LinksController{
		BaseController: &controller.BaseController{},
	}
	linksHandlerFunc := handler.MakeHandler(handler.Handler, &linksController, fileSystem, db)

	skillReviewsController := controller.SkillReviewsController{
		BaseController: &controller.BaseController{},
	}
	skillReviewsHandlerFunc := handler.MakeHandler(handler.Handler, &skillReviewsController, fileSystem, db)

	skillIconsController := controller.SkillIconsController{
		BaseController: &controller.BaseController{},
	}
	skillIconsHandlerFunc := handler.MakeHandler(handler.Handler, &skillIconsController, fileSystem, db)

	usersController := controller.UsersController{
		BaseController: &controller.BaseController{},
	}
	usersHandlerFunc := handler.MakeHandler(handler.Handler, &usersController, fileSystem, db)

	routes = []Route{
		{"/api/skills/", skillsHandlerFunc},
		{"/api/skills", skillsHandlerFunc},
		{"/api/teammembers/", teamMembersHandlerFunc},
		{"/api/teammembers", teamMembersHandlerFunc},
		{"/api/tmskills/", tmSkillsHandlerFunc},
		{"/api/tmskills", tmSkillsHandlerFunc},
		{"/api/links/", linksHandlerFunc},
		{"/api/links", linksHandlerFunc},
		{"/api/skillreviews", skillReviewsHandlerFunc},
		{"/api/skillreviews/", skillReviewsHandlerFunc},
		{"/api/skillicons", skillIconsHandlerFunc},
		{"/api/skillicons/", skillIconsHandlerFunc},
		{"/api/users", usersHandlerFunc},
		{"/api/users/", usersHandlerFunc},
	}
}

/*
StartRouter() instantiates a new http.ServeMux and registers with it each
endpoint that is currently being handled by the SkillDirectory REST API with an
appropriate handler function for that endpoint. This http.ServeMux is returned.
*/
func StartRouter() (mux *http.ServeMux) {
	initCassandra()
	initFileSystem()
	loadRoutes()
	mux = http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handlerFunc)
	}
	return
}
