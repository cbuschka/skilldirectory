package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/model"
	"testing"
	"skilldirectory/data"
)

func TestTeamMemberControllerBase(t *testing.T) {
	base := BaseController{}
	tc := TeamMembersController{BaseController: &base}

	if base != *tc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllTeamMembers(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/teammembers", nil)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllTeamMembers_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/teammembers", nil)
	tc := getTeamMembersController(request, &data.MockErrorDataAccessor{})

	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetTeamMember(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/teammembers/1234", nil)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTeamMember_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/teammebers/1234", nil)
	tc := getTeamMembersController(request, &data.MockErrorDataAccessor{})
	
	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTeamMember(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/teammembers/1234", nil)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTeamMember_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/teammembers/1234", nil)
	tc := getTeamMembersController(request, &data.MockErrorDataAccessor{})

	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTeamMember_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/teammembers/", nil)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostTeamMember(t *testing.T) {
	body := getReaderForNewTeamMember("1234", "Joe Smith", "Cabbage Plucker")
	request := httptest.NewRequest(http.MethodPost, "/teammembers", body)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostTeamMember_NoName(t *testing.T) {
	body := getReaderForNewTeamMember("1234", "", "Cabbage Plucker")
	request := httptest.NewRequest(http.MethodPost, "/teammembers", body)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in TeamMember POST request.", "name")
	}
}

func TestPostTeamMember_NoTitle(t *testing.T) {
	body := getReaderForNewTeamMember("1234", "Joe Smith", "")
	request := httptest.NewRequest(http.MethodPost, "/teammembers", body)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in TeamMember POST request.", "title")
	}
}

func TestPostTeamMember_NoTeamMember(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/teammembers", nil)
	tc := getTeamMembersController(request, &data.MockDataAccessor{})
	
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostTeamMember_Error(t *testing.T) {
	body := getReaderForNewTeamMember("1234", "Joe Smith", "Cabbage Plucker")
	request := httptest.NewRequest(http.MethodPost, "/teammembers", body)
	tc := getTeamMembersController(request, &data.MockErrorDataAccessor{})
	
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

/*
getTeamMembersController is a helper function for creating and initializing a new BaseController with
the given HTTP request and DataAccessor. Returns a new TeamMembersController created with that BaseController.
 */
func getTeamMembersController(request *http.Request, dataAccessor data.DataAccess) TeamMembersController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor)
	return TeamMembersController{BaseController: &base}
}

/*
getReaderForNewTeamMember is a helper function for a new TeamMember with the given id, name, and title.
This TeamMember is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
 */
func getReaderForNewTeamMember(id, name, title string) *bytes.Reader {
	newTeamMember := model.NewTeamMember(id, name, title)
	b, _ := json.Marshal(newTeamMember)
	return bytes.NewReader(b)
}
