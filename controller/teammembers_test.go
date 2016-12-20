package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/model"
	"testing"
)

func TestTeamMembersControllerBase(t *testing.T) {
	base := BaseController{}
	tc := TeamMembersController{BaseController: &base}
	if base != *tc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllTeamMembers(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/teammembers", nil), &MockDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllTeamMembersError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/teammembers", nil), &MockErrorDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetTeamMember(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/teammembers/1234", nil), &MockDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTeamMemberError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/teammembers/1234", nil), &MockErrorDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTeamMember(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/teammembers/1234", nil), &MockDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTeamMemberError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/teammembers/1234", nil), &MockErrorDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTeamMemberNoKey(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/teammembers/", nil), &MockDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostTeamMember(t *testing.T) {
	base := BaseController{}
	newTM := model.NewTeamMember("1234", "Joe Smith", "Cabbage Plucker")
	b, _ := json.Marshal(newTM)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/teammembers", reader), &MockDataAccessor{})

	tc := TeamMembersController{BaseController: &base}
	err := tc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostTeamMemberNoName(t *testing.T) {
	base := BaseController{}
	newTM := model.NewTeamMember("1234", "", "Cabbage Plucker")
	b, _ := json.Marshal(newTM)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/teammembers", reader), &MockDataAccessor{})

	tc := TeamMembersController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to not specifying value for \"Name\" field in TeamMember POST request.")
	}
}

func TestPostTeamMemberNoTitle(t *testing.T) {
	base := BaseController{}
	newTM := model.NewTeamMember("1234", "Joe Smith", "")
	b, _ := json.Marshal(newTM)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/teammembers", reader), &MockDataAccessor{})

	tc := TeamMembersController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to not specifying value for \"Title\" field in TeamMember POST request.")
	}
}

func TestPostNoTeamMember(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/teammembers/", nil), &MockDataAccessor{})
	tc := TeamMembersController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostTeamMemberError(t *testing.T) {
	base := BaseController{}
	newTM := model.NewTeamMember("1234", "Joe Smith", "Cabbage Plucker")
	b, _ := json.Marshal(newTM)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/teammembers", reader), &MockErrorDataAccessor{})

	tc := TeamMembersController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}
