package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/model"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestTMSkillsController_Base(t *testing.T) {
	base := BaseController{}
	tc := TMSkillsController{BaseController: &base}

	if base != *tc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllTMSkills(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/tmskills", nil)
	tc := getTMSkillsController(request, false)

	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllTMSkills_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/tmskills", nil)
	tc := getTMSkillsController(request, true)

	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetTMSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/tmskills/1234", nil)
	tc := getTMSkillsController(request, false)

	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTMSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/tmskills/1234", nil)
	tc := getTMSkillsController(request, true)

	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/tmskills/1234", nil)
	tc := getTMSkillsController(request, false)

	err := tc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/tmskills/1234", nil)
	tc := getTMSkillsController(request, true)

	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkill_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/tmskills/", nil)
	tc := getTMSkillsController(request, false)

	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostTMSkill(t *testing.T) {
	body := getReaderForNewTMSkill(1234, 2345, 3456)
	request := httptest.NewRequest(http.MethodPost, "/api/tmskills", body)
	tc := getTMSkillsController(request, false)

	err := tc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostTMSkill_NoSkillID(t *testing.T) {
	body := getReaderForNewTMSkill(1234, 0, 3456)
	request := httptest.NewRequest(http.MethodPost, "/api/tmskills", body)
	tc := getTMSkillsController(request, false)

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in TMSkill POST"+
			" request.", "skill_id")
	}
}

func TestPostTMSkill_NoTeamMemberID(t *testing.T) {
	body := getReaderForNewTMSkill(1234, 2345, 0)
	request := httptest.NewRequest(http.MethodPost, "/api/tmskills", body)
	tc := getTMSkillsController(request, false)

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in TMSkill POST"+
			" request.", "team_member_id")
	}
}

func TestPostTMSkill_NoTMSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/tmskills", nil)
	tc := getTMSkillsController(request, false)

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostTMSkill_Error(t *testing.T) {
	body := getReaderForNewTMSkill(1234, 2345, 3456)
	request := httptest.NewRequest(http.MethodPost, "/api/tmskills", body)
	tc := getTMSkillsController(request, true)

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func Test_validateTMSkillFields(t *testing.T) {
	tc := getTMSkillsController(nil, true)
	tmSkill := model.TMSkill{
		SkillID: 1234,
	}
	err := tc.validateTMSkillFields(tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect empty " +
			"TMSkill.TeamMemberID field.")
	}

	tmSkill = model.TMSkill{
		TeamMemberID: 1234,
	}
	err = tc.validateTMSkillFields(tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect empty " +
			"TMSkill.SkillID field.")
	}

	tmSkill = model.TMSkill{
		SkillID:      1234,
		TeamMemberID: 1234,
	}
	err = tc.validateTMSkillFields(tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect invalid " +
			"ID field.")
	}

	tmSkill = model.TMSkill{
		SkillID:      1234,
		TeamMemberID: 1234,
		Proficiency:  9000,
	}
	err = tc.validateTMSkillFields(tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect invalid " +
			"TMSkill.Proficiency field.")
	}
}

func TestUpdateTMSkill(t *testing.T) {
	body := getReaderForNewTMSkill(1234, 2345, 3456)
	request := httptest.NewRequest(http.MethodPut, "/api/tmskills/1234", body)
	tc := getTMSkillsController(request, false)

	err := tc.Put()
	if err != nil {
		t.Errorf("Put failed: %s", err.Error())
	}
}

func TestUpdateTMSkillError(t *testing.T) {
	body := getReaderForNewTMSkill(1234, 2345, 3456)
	request := httptest.NewRequest(http.MethodPut, "/api/tmskills/1234", body)
	tc := getTMSkillsController(request, true)

	err := tc.Put()
	if err == nil {
		t.Errorf("Expecting error on update")
	}
}

func TestUpdateTMSkillNoID(t *testing.T) {
	body := getReaderForNewTMSkill(1234, 2345, 3456)
	request := httptest.NewRequest(http.MethodPost, "/api/tmskills", body)
	tc := getTMSkillsController(request, false)

	err := tc.Put()
	if err == nil {
		t.Errorf("Expecting error on update with no id in url")
	}
}

func TestValidProf(t *testing.T) {
	tmSkill := model.NewTMSkillSetDefaults(1, 2, 3, 0)
	c := getTMSkillsController(nil, false)
	err := c.validateTMSkillFields(tmSkill)
	if err != nil {
		t.Errorf("Expecting a valid tmskill: %v", tmSkill)
	}
}

func TestTMSkillOptions(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/tmskills", nil)
	tsc := getTMSkillsController(request, false)

	err := tsc.Options()
	if err != nil {
		t.Errorf("OPTIONS requests should always return a 200 response.")
	}
	if tsc.w.Header().Get("Access-Control-Allow-Methods") != "PUT, "+GetDefaultMethods() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Methods' contains" +
			" incorrect value")
	}
	if tsc.w.Header().Get("Access-Control-Allow-Headers") != GetDefaultHeaders() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Headers' contains" +
			" incorrect value")
	}
}

/*
getTMSkillsController is a helper function for creating and initializing a new
BaseController with the given HTTP request and DataAccessor. Returns a new
TMSkillsController created with that BaseController.
*/
func getTMSkillsController(request *http.Request, errSwitch bool) TMSkillsController {
	base := BaseController{}
	base.SetTest(errSwitch)
	base.InitWithGorm(httptest.NewRecorder(), request, nil, logrus.New(), nil)
	return TMSkillsController{BaseController: &base}
}

/*
getReaderForNewTMSkill is a helper function for a new TMSkill with the given id,
skillID, and teamMemberID. This TMSkill is then marshaled into JSON. A new Reader
is created and returned for the resulting []byte.
*/
func getReaderForNewTMSkill(id, skillID, teamMemberID uint) *bytes.Reader {
	newTMSkill := model.NewTMSkillDefaults(id, skillID, teamMemberID)
	b, _ := json.Marshal(newTMSkill)
	return bytes.NewReader(b)
}
