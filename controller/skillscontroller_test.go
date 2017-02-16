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

func TestSkillsControllerBase(t *testing.T) {
	base := BaseController{}
	sc := SkillsController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expecting Base() to return base pointer")
	}
}

func TestGetAllSkills(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skills", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllSkills_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skills", nil)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skills/1234", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skills/1234", nil)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skills/1234", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Delete()
	if err != nil {
		t.Errorf("Expected no error, but got one: %s", err.Error())
	}
}

func TestDeleteSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skills/1234", nil)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkill_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skills/", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostSkill(t *testing.T) {
	body := getReaderForNewSkill("", "BestSkillNameEver", model.ScriptedSkillType)
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/api/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkill_NoName(t *testing.T) {
	body := getReaderForNewSkill("1234", "", model.ScriptedSkillType)
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/api/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Skill POST request.", "name")
	}
}

func TestPostSkill_NoSkillType(t *testing.T) {
	body := getReaderForNewSkill("1234", "SomeName", "")
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/api/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Skill POST request.", "skill_type")
	}
}

func TestPostSkill_InvalidType(t *testing.T) {
	body := getReaderForNewSkill("", "", "badtype")
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/api/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkill_NoSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/skills/", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkill_Error(t *testing.T) {
	body := getReaderForNewSkill("", "", model.ScriptedSkillType)
	request := httptest.NewRequest(http.MethodPost, "/api/skills", body)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

/*
getSkillsController is a helper function for creating and initializing a new BaseController with
the given HTTP request and DataAccessor. Returns a new SkillsController created with that BaseController.
*/
func getSkillsController(request *http.Request, dataAccessor data.DataAccess) SkillsController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor, nil, logrus.New())
	return SkillsController{BaseController: &base}
}

/*
getReaderForNewSkill is a helper function for a new Skill with the given id, name, and skillType.
This Skill is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewSkill(id, name, skillType string) *bytes.Reader {
	newSkill := model.NewSkill(id, name, skillType)
	b, _ := json.Marshal(newSkill)
	return bytes.NewReader(b)
}
