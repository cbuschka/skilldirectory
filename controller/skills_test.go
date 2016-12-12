package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"skilldirectory/data"
	"testing"
)

type MockDataAccessor struct{}

func (m MockDataAccessor) Save(s string, i interface{}) error { return nil }
func (m MockDataAccessor) Read(s string, i interface{}) error { return nil }
func (m MockDataAccessor) Delete(s string) error              { return nil }
func (m MockDataAccessor) ReadAll(s string, r data.ReadAllInterface) ([]interface{}, error) {
	return nil, nil
}

type MockErrorDataAccessor struct{}

func (e MockErrorDataAccessor) Save(s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Read(s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Delete(s string) error              { return fmt.Errorf("") }
func (e MockErrorDataAccessor) ReadAll(s string, r data.ReadAllInterface) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}

func TestBase(t *testing.T) {
	base := BaseController{}
	sc := SkillsController{BaseController: &base}
	if base != *sc.Base() {
		t.Error("Expecting Base() to return base pointer")
	}
}

func TestGetAll(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/skills", nil), &MockDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/skills", nil), &MockErrorDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkill(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/skills/1234", nil), &MockDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSkillError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/skills/1234", nil), &MockErrorDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDelete(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/skills/1234", nil), &MockDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteError(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/skills/1234", nil), &MockErrorDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteNoKey(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/skills/", nil), &MockDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}
