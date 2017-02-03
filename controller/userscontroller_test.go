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
	tc := UsersController{BaseController: &base}

	if base != *tc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestPostUser(t *testing.T) {
	body := getReaderForUserLogin("test", "test")
	request := httptest.NewRequest(http.MethodPost, "/users", body)
	tc := getUsersController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostUser_NoLogin(t *testing.T) {
	body := getReaderForUserLogin("", "test")
	request := httptest.NewRequest(http.MethodPost, "/users", body)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in User POST request.", "name")
	}
}

func TestPostUser_NoPassword(t *testing.T) {
	body := getReaderForUserLogin("test", "")
	request := httptest.NewRequest(http.MethodPost, "/users", body)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in User POST request.", "title")
	}
}

func TestPostUser_Error(t *testing.T) {
	body := getReaderForUserLogin("test", "test")
	request := httptest.NewRequest(http.MethodPost, "/users", body)
	tc := getTeamMembersController(request, &data.MockErrorDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
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

/*
getReaderForNewUser is a helper function for a new User with the given login and password.
This User is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForUserLogin(login, password string) *bytes.Reader {
	newUser := model.NewUser(login, password)
	b, _ := json.Marshal(newUser)
	return bytes.NewReader(b)
}
