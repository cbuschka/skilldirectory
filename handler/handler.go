package handler

import (
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
	"skilldirectory/errors"
	"skilldirectory/util"
	"sync"

	"github.com/jinzhu/gorm"
)

// This mutex will prevent race conditions on concurrent requests
var mutex = &sync.Mutex{}

/*
MakeHandler() returns a new function of the adapter type http.HandlerFunc using
the passed-in function, fn.
*/
func MakeHandler(
	fn func(http.ResponseWriter, *http.Request, controller.RESTController,
		data.FileSystem, *gorm.DB),
	cont controller.RESTController, fs data.FileSystem, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r, cont, fs, db)
	}
}

/*
Handler() should be invoked to handle responding to the passed-in HTTP request.
Responses are sent via the passed-in http.ResponseWriter.

The passed-in RESTController is first initialized using the specified
http.ResponseWriter and http.Request, and is connected to the Skills database.
Once initialized, the RESTController is used to handle responses to the
passed-in HTTP request.

If the RESTController generates any errors, then Handler() will
log them, and respond to the request with the appropriate error.
*/
func Handler(w http.ResponseWriter, r *http.Request, cont controller.RESTController, fs data.FileSystem, db *gorm.DB) {

	//Lock the critical section
	mutex.Lock()
	//Make sure we always unlock when the function extends
	defer mutex.Unlock()

	log := util.LogInit()
	log.Printf("Handling Request: [%s] Path: [%s]", r.Method, r.RequestURI)
	log.Debugf("Request: %s", r.Body)
	cont.Base().InitWithGorm(w, r, fs, log, db)

	var err error
	switch r.Method {
	case http.MethodGet:
		err = cont.Get()
	case http.MethodPost:
		err = cont.Post()
	case http.MethodDelete:
		err = cont.Delete()
	case http.MethodPut:
		err = cont.Put()
	case http.MethodOptions:
		err = cont.Options()
	}

	var statusCode int
	if err != nil {
		switch err.(type) {
		case errors.MarshalingError, errors.InvalidSkillTypeError,
			errors.MissingIDError, errors.IncompletePOSTBodyError,
			errors.InvalidPOSTBodyError, errors.InvalidPUTBodyError:
			statusCode = http.StatusBadRequest
		case errors.SavingError, errors.ReadError:
			statusCode = http.StatusInternalServerError
		case errors.NoSuchIDError:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		log.Warnf("Handler Method: %s, %T: %v", r.Method, err, err)
		http.Error(w, err.Error(), statusCode)
	}
}
