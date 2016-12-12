package controller

import (
	"fmt"
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

// func TestGet(t *testing.T) {
// 	fmt.Println("testing")
// 	base := BaseController{}
// 	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "skills", nil), &MockDataAccessor{})
// 	fmt.Println("testing")
//
// 	sc := SkillsController{BaseController: &base}
// 	err := sc.Get()
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// }
