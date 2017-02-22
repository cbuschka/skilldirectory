package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/gormmodel"
	"skilldirectory/model"
	"testing"

	"github.com/Sirupsen/logrus"
)

//
func TestSkillsControllerBase(t *testing.T) {
	base := BaseController{}
	sc := SkillsController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expecting Base() to return base pointer")
	}
}

func TestGetAllSkills(t *testing.T) {
	sc := getSkillsController(
		httptest.NewRequest(
			http.MethodGet,
			"/api/skills",
			nil),
		false)

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllSkills_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skills", nil)
	sc := getSkillsController(request, true)

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skills/1234", nil)
	sc := getSkillsController(request, false)

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skills/1234", nil)
	sc := getSkillsController(request, true)

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skills/1234", nil)
	sc := getSkillsController(request, false)

	err := sc.Delete()
	if err != nil {
		t.Errorf("Expected no error, but got one: %s", err.Error())
	}
}

func TestDeleteSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skills/1234", nil)
	sc := getSkillsController(request, true)

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestDeleteSkill_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skills/", nil)
	sc := getSkillsController(request, false)

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key")
	}
}

func TestDeleteSkill_0Key(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skills/0", nil)
	sc := getSkillsController(request, false)

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when key is 0")
	}
}

func TestPostSkill(t *testing.T) {
	sc := getSkillsController(
		httptest.NewRequest(
			http.MethodPost,
			"/api/skills",
			getReaderForNewSkill(0, "BestSkillNameEver", model.ScriptedSkillType)),
		false)

	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkill_NoName(t *testing.T) {
	sc := getSkillsController(
		httptest.NewRequest(
			http.MethodPost,
			"/api/skills",
			getReaderForNewSkill(1234, "", model.ScriptedSkillType)),
		false)

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Skill POST request.", "name")
	}
}

func TestPostSkill_NoSkillType(t *testing.T) {
	sc := getSkillsController(
		httptest.NewRequest(
			http.MethodPost,
			"/api/skills", getReaderForNewSkill(1234, "SomeName", "")),
		false)

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Skill POST request.", "skill_type")
	}
}

func TestPostSkill_InvalidType(t *testing.T) {
	sc := getSkillsController(
		httptest.NewRequest(
			http.MethodPost,
			"/api/skills", getReaderForNewSkill(0, "", "badtype")),
		false)

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkill_NoSkill(t *testing.T) {
	sc := getSkillsController(
		httptest.NewRequest(http.MethodPost, "/api/skills/", nil),
		false)

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkill_Error(t *testing.T) {
	sc := getSkillsController(
		httptest.NewRequest(
			http.MethodPost, "/api/skills",
			getReaderForNewSkill(0, "", model.ScriptedSkillType)),
		true)

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

//
func TestSkillOptions(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/skills", nil)
	sc := getSkillsController(request, false)

	err := sc.Options()
	if err != nil {
		t.Errorf("OPTIONS requests should always return a 200 response.")
	}
	if sc.w.Header().Get("Access-Control-Allow-Methods") != GetDefaultMethods() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Methods' contains" +
			" incorrect value")
	}
	if sc.w.Header().Get("Access-Control-Allow-Headers") != GetDefaultHeaders() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Headers' contains" +
			" incorrect value")
	}
}

//
// /*
// getSkillsController is a helper function for creating and initializing a new BaseController with
// the given HTTP request and DataAccessor. Returns a new SkillsController created with that BaseController.
// */
func getSkillsController(request *http.Request, errSwitch bool) SkillsController {
	base := BaseController{}
	base.SetTest(errSwitch)
	base.Init(httptest.NewRecorder(), request, nil, nil, logrus.New())
	return SkillsController{BaseController: &base}
}

/*
getReaderForNewSkill is a helper function for a new Skill with the given id, name, and skillType.
This Skill is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewSkill(id uint, name, skillType string) *bytes.Reader {
	newSkill := gormmodel.NewSkill(id, name, skillType)
	b, _ := json.Marshal(newSkill)
	return bytes.NewReader(b)
}
