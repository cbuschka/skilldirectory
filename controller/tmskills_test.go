package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/model"
	"testing"
)

func TestTMSkillsController_Base(t *testing.T) {
	base := BaseController{}
	tc := TMSkillsController{BaseController: &base}
	if base != *tc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllTMSkills(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/tmskills", nil), &MockDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllTMSkillsError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/tmskills", nil), &MockErrorDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetTMSkill(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/tmskills/1234", nil), &MockDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTMSkillError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/tmskills/1234", nil), &MockErrorDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkill(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/tmskills/1234", nil), &MockDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkillError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/tmskills/1234", nil), &MockErrorDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkillNoKey(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/tmskills/", nil), &MockDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostTMSkill(t *testing.T) {
	base := BaseController{}
	newTMSkill := model.NewTMSkillDefaults("1234", "2345", "3456")
	b, _ := json.Marshal(newTMSkill)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/tmskills", reader), &MockDataAccessor{})

	tc := TMSkillsController{BaseController: &base}
	err := tc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostTMSkillNoSkillID(t *testing.T) {
	base := BaseController{}
	newTMSkill := model.NewTMSkillDefaults("1234", "", "3456")
	b, _ := json.Marshal(newTMSkill)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/tmskills", reader), &MockDataAccessor{})

	tc := TMSkillsController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to not specifying value for \"name\" field in TMSkill POST request.")
	}
}

func TestPostTMSKillNoTeamMemberID(t *testing.T) {
	base := BaseController{}
	newTMSkill := model.NewTMSkillDefaults("1234", "2345", "")
	b, _ := json.Marshal(newTMSkill)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/tmskills", reader), &MockDataAccessor{})

	tc := TMSkillsController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to not specifying value for \"title\" field in TMSkill POST request.")
	}
}

func TestPostNoTMSkill(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/tmskills/", nil), &MockDataAccessor{})
	tc := TMSkillsController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostTMSkillError(t *testing.T) {
	base := BaseController{}
	newTMSkill := model.NewTMSkillDefaults("1234", "2345", "3456")
	b, _ := json.Marshal(newTMSkill)
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/tmskills", reader), &MockErrorDataAccessor{})

	tc := TMSkillsController{BaseController: &base}
	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}
