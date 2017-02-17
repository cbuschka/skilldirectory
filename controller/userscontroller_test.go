package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/data"
	"skilldirectory/model"
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
	request := httptest.NewRequest(http.MethodGet, "/api/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDelete(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockErrorDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPut(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/api/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockErrorDataAccessor{})

	err := sc.Put()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestOptions(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/users", bytes.NewBufferString(""))
	sc := getUsersController(request, &data.MockDataAccessor{})

	err := sc.Options()
	if err != nil {
		t.Error(err)
	}
}

func TestPostUser_NoCode(t *testing.T) {
	body := getReaderForNewCredentials("", "")
	request := httptest.NewRequest(http.MethodPost, "/api/users", body)
	uc := getUsersController(request, &data.MockDataAccessor{})

	err := uc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostUser_NoClientID(t *testing.T) {
	body := getReaderForNewCredentials("foobarbaz", "")
	request := httptest.NewRequest(http.MethodPost, "/api/users", body)
	uc := getUsersController(request, &data.MockDataAccessor{})

	err := uc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

/*
getReaderForNewUser is a helper function for a new Skill with the given id, name, and skillType.
This Skill is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewCredentials(code, clientID string) *bytes.Reader {
	newCredentials := model.AuthCredentials{Code: code, Id: clientID}
	b, _ := json.Marshal(newCredentials)
	return bytes.NewReader(b)
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
