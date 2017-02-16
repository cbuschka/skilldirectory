package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"skilldirectory/data"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestUsersControllerBase(t *testing.T) {
	base := BaseController{}
	sc := UsersController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expecting Base() to return base pointer")
	}
}

func TestGet(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDelete(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockErrorDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPut(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockErrorDataAccessor{})

	err := sc.Put()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestOptions(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockDataAccessor{})

	err := sc.Options()
	if err != nil {
		t.Error(err)
	}
}

/*
getUsersController is a helper function for creating and initializing a new BaseController with
the given HTTP request and DataAccessor. Returns a new UsersController created with that BaseController.
*/
func getUsersController(request *http.Request, dataAccessor data.DataAccess) UsersController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor, nil, logrus.New())
	return UsersController{BaseController: &base}
}
