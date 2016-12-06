package handler

import (
  "testing"
  "fmt"
  "skilldirectory/data"
  "net/http"
  "net/http/httptest"
  "skilldirectory/model"
  "encoding/json"
  "bytes"



)

type MockDataAccessor struct {}
func (m MockDataAccessor)Save(s string, i interface{}) error {return nil}
func (m MockDataAccessor) Read(s string, i interface{}) error {return nil}
func (m MockDataAccessor) Delete(s string) error {return nil}

type MockErrorDataAccessor struct {}
func (e MockErrorDataAccessor)Save(s string, i interface{}) error {return fmt.Errorf("")}
func (e MockErrorDataAccessor) Read(s string, i interface{}) error {return fmt.Errorf("")}
func (e MockErrorDataAccessor) Delete(s string) error {return fmt.Errorf("")}

func TestLoadSkill(t *testing.T) {
  fmt.Println("PASS")
  skillsConnector = data.NewAccessor(MockDataAccessor{})
  _, err := loadSkill("1234")
  if err != nil {
    t.Errorf("Load skill error %s", err.Error())
  }
}

func TestLoadSkillError(t *testing.T) {
  fmt.Println("PASS")
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
  SkillsHandler(w,r, "")
  if w.Result().StatusCode != 200 {
    t.Errorf("Expected 200 from response")
  }
  if w.Body == nil {
    t.Errorf("Expected reponse body")
  }
}

func TestSkillsHandlerPOSTSuccess(t *testing.T) {
  skillsConnector = data.NewAccessor(MockDataAccessor{})
  b,_ := json.Marshal(model.NewSkill("", "", "scripted"))
  reader := bytes.NewReader(b)
  r := httptest.NewRequest(http.MethodPost, "/", reader)
  w := httptest.NewRecorder()
  SkillsHandler(w,r, "/")
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
  SkillsHandler(w,r, "")
  if w.Result().StatusCode != 404 {
    t.Errorf("Expected 404 from bad request")
  }
}
