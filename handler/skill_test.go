package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"skilldirectory/data"
	"skilldirectory/model"
	"testing"
)

type MockDataAccessor struct{}

func (m MockDataAccessor) Save(s string, i interface{}) error { return nil }
func (m MockDataAccessor) Read(s string, i interface{}) error { return nil }
func (m MockDataAccessor) Delete(s string) error              { return nil }
func (m MockDataAccessor) ReadAll(s string, r data.ReadAllInterface) ([]interface{}, error) {
	return nil, nil
}
func (d MockDataAccessor) FilteredReadAll(s string, r data.ReadAllInterface, f func (interface{}) bool) ([]interface{}, error) {
	return nil, nil
}

type MockErrorDataAccessor struct{}

func (e MockErrorDataAccessor) Save(s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Read(s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Delete(s string) error              { return fmt.Errorf("") }
func (e MockErrorDataAccessor) ReadAll(s string, r data.ReadAllInterface) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}
func (d MockErrorDataAccessor) FilteredReadAll(s string, r data.ReadAllInterface, f func (interface{}) bool) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}

func TestLoadSkill(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	_, err := loadSkill("1234")
	if err != nil {
		t.Errorf("Load skill error %s", err.Error())
	}
}

func TestLoadSkillError(t *testing.T) {
	skillsConnector = data.NewAccessor(MockErrorDataAccessor{})
	_, err := loadSkill("1234")
	if err == nil {
		t.Errorf("Expecting load skill error, got nil")
	}
}

func TestSkillsHandlerGETSuccess(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	r := httptest.NewRequest(http.MethodGet, "/1234", nil)
	w := httptest.NewRecorder()
	SkillsHandler(w, r, "")
	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 from response")
	}
	if w.Body == nil {
		t.Errorf("Expected reponse body")
	}
}

func TestSkillsHandlerPOSTSuccess(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	b, _ := json.Marshal(model.NewSkill("", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	r := httptest.NewRequest(http.MethodPost, "/", reader)
	w := httptest.NewRecorder()
	SkillsHandler(w, r, "/")
	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 from response")
	}
	if w.Body == nil {
		t.Errorf("Expected reponse body")
	}
}

func TestSkillsHandlerError(t *testing.T) {
	skillsConnector = data.NewAccessor(MockErrorDataAccessor{})
	r := httptest.NewRequest(http.MethodGet, "/1234", nil)
	w := httptest.NewRecorder()
	SkillsHandler(w, r, "")
	if w.Result().StatusCode != 404 {
		t.Errorf("Expected 404 from bad request")
	}
}

func TestInvalidSkill(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	b, _ := json.Marshal(model.NewSkill("", "", "BadSkillType"))
	reader := bytes.NewReader(b)
	r := httptest.NewRequest(http.MethodPost, "/", reader)
	err := addSkill(r)
	if err == nil {
		t.Errorf("Expecting an error for BadSkillType")
	}
}

func TestNilSkill(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	err := addSkill(r)
	if err == nil {
		t.Errorf("Expecting an error for nil body")
	}
}

func TestAddSkillError(t *testing.T) {
	skillsConnector = data.NewAccessor(MockErrorDataAccessor{})
	b, _ := json.Marshal(model.NewSkill("", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	r := httptest.NewRequest(http.MethodPost, "/", reader)
	err := addSkill(r)
	if err == nil {
		t.Errorf("Expecting an error for BadSkillType")
	}
}

func TestRemoveSkill(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	b, _ := json.Marshal(model.NewSkill("1234", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	r := httptest.NewRequest(http.MethodDelete, "/skills/1234", reader)
	err := removeSkill(r)
	if err != nil {
		t.Errorf("Did not expect error when deleting skill")
	}
}

func TestRemoveSkillError(t *testing.T) {
	skillsConnector = data.NewAccessor(MockErrorDataAccessor{})
	b, _ := json.Marshal(model.NewSkill("", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	r := httptest.NewRequest(http.MethodDelete, "/skills", reader)
	err := removeSkill(r)
	if err == nil {
		t.Errorf("Expected error due to skill deletion failure")
	}
}

func TestPerformGet(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	r := httptest.NewRequest(http.MethodGet, "/skills", nil)

	w := httptest.NewRecorder()
	performGet(w, r)
}

func TestPerformGetError(t *testing.T) {
	skillsConnector = data.NewAccessor(MockErrorDataAccessor{})
	r := httptest.NewRequest(http.MethodGet, "/skills", nil)

	w := httptest.NewRecorder()
	err := performGet(w, r)
	if err == nil {
		t.Errorf("Expecting error for TestPerform")
	}
}

func TestUrlOptions(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	urls := []string{"/skills", "/skills/", "/skills?skilltype=scripted", "/skills/sdfasdfas"}
	for _, uri := range urls {
		r := httptest.NewRequest(http.MethodGet, uri, nil)
		w := httptest.NewRecorder()
		err := performGet(w, r)
		if err != nil {
			t.Errorf("Perform Get Error From URL: %s, %s", uri, err.Error())
		}
		if w.Code != 200 {
			t.Errorf("Perform Get Error From URL: %s, %d", uri, w.Code)
		}
	}

}
