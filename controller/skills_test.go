package controller

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
func (d MockDataAccessor) FilteredReadAll(s string, r data.ReadAllInterface, f func(interface{}) bool) ([]interface{}, error) {
	return nil, nil
}

type MockErrorDataAccessor struct{}

func (e MockErrorDataAccessor) Save(s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Read(s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Delete(s string) error              { return fmt.Errorf("") }
func (e MockErrorDataAccessor) ReadAll(s string, r data.ReadAllInterface) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}
func (d MockErrorDataAccessor) FilteredReadAll(s string, r data.ReadAllInterface, f func(interface{}) bool) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////// MockInMemoryDataAccessor /////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
MockInMemoryDataAccessor implements the data.dataccessor.DataAccess interface by providing in-memory data storage,
designed to facilitate easy testing of components that rely on skill storage and retrieval without requiring access to
an external filesystem or database when running tests. Because MockInMemoryDataAccessor uses a computer's volatile
memory, it should not be used for permanent data storage, and is unlikely to have a use outside of unit testing.
*/
type MockInMemoryDataAccessor struct {
	dataMap map[string][]byte // map of object to byte slice
}

func NewMockInMemoryDataAccessor() MockInMemoryDataAccessor {
	var retVal MockInMemoryDataAccessor
	retVal.dataMap = make(map[string][]byte)
	return retVal
}

func (e MockInMemoryDataAccessor) Save(ID string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	e.dataMap[ID] = b
	return nil
}

func (e MockInMemoryDataAccessor) Read(ID string, object interface{}) error {
	data := e.dataMap[ID]
	if len(data) == 0 {
		return fmt.Errorf("No such object with ID: %s", ID)
	}
	json.Unmarshal(data, &object)
	return nil
}

func (e MockInMemoryDataAccessor) Delete(ID string) error {
	fmt.Println("Deleting this key from map:", ID)
	data := e.dataMap[ID]
	if len(data) == 0 {
		return fmt.Errorf("No such object with ID: %s", ID)
	}
	e.dataMap[ID] = make([]byte, 0)
	return nil
}

func (e MockInMemoryDataAccessor) ReadAll(path string, readType data.ReadAllInterface) ([]interface{}, error) {
	returnObjects := []interface{}{}
	object := readType.GetType()
	for _, val := range e.dataMap {
		json.Unmarshal(val, object)
		returnObjects = append(returnObjects, object)
	}
	return returnObjects, nil
}

func (e MockInMemoryDataAccessor) FilteredReadAll(path string, readType data.ReadAllInterface,
	filterFunc func(interface{}) bool) ([]interface{}, error) {
	returnObjects := []interface{}{}
	object := readType.GetType()
	for _, val := range e.dataMap {
		json.Unmarshal(val, &object)
		if filterFunc(object) {
			returnObjects = append(returnObjects, object)
		}
	}
	return returnObjects, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

func TestGetSkillsFiltered(t *testing.T) {
	base := BaseController{}
	skillsConnector := SkillsController{BaseController: &base}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/skills?skilltype=scripted", nil)
	accessor := NewMockInMemoryDataAccessor()
	base.Init(responseRecorder, request, &accessor)

	newScriptedSkill := model.NewSkill("1234", "TestSkillName", model.ScriptedSkillType)
	skillsConnector.session.Save(newScriptedSkill.ID, newScriptedSkill)
	newCompiledSkill := model.NewSkill("2136", "TestSkillName", model.CompiledSkillType)
	skillsConnector.session.Save(newCompiledSkill.ID, newCompiledSkill)

	err := skillsConnector.Get()
	if err != nil {
		t.Errorf("Did not expect error when getting skills with filter")
	}

	correctResponseBody := "[{\"ID\":\"1234\",\"Name\":\"TestSkillName\",\"SkillType\":\"scripted\"}]"
	if responseRecorder.Body.String() != correctResponseBody {
		t.Errorf("Failed to properly filter based on skilltype. "+
			"Expected Response body to be \n\t %s\n But got\n\t %s\n",
			correctResponseBody, responseRecorder.Body.String())
	}
}

func TestGetSkillsFilteredBadSkillType(t *testing.T) {
	base := BaseController{}
	skillsConnector := SkillsController{BaseController: &base}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/skills?skilltype=invalid-skill-type", nil)
	accessor := NewMockInMemoryDataAccessor()
	base.Init(responseRecorder, request, &accessor)

	err := skillsConnector.Get()
	if err == nil {
		t.Errorf("Expected error due to invalid skill type.")
	}
}

func TestGetSkillsFilteredError(t *testing.T) {
	base := BaseController{}
	skillsConnector := SkillsController{BaseController: &base}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/skills?skilltype=scripted", nil)
	accessor := MockErrorDataAccessor{}
	base.Init(responseRecorder, request, &accessor)

	err := skillsConnector.Get()
	if err == nil {
		t.Errorf("Expecting error for TestGetSkillsFilteredError")
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

func TestPostSkill(t *testing.T) {
	base := BaseController{}
	b, _ := json.Marshal(model.NewSkill("1234", "SomeName", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/skills", reader), &MockDataAccessor{})

	sc := SkillsController{BaseController: &base}
	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkillNoName(t *testing.T) {
	base := BaseController{}
	b, _ := json.Marshal(model.NewSkill("1234", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/skills", reader), &MockDataAccessor{})

	sc := SkillsController{BaseController: &base}
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to not specifying value for \"Name\" field in Skill POST request.")
	}
}

func TestPostSkillNoSkillType(t *testing.T) {
	base := BaseController{}
	b, _ := json.Marshal(model.NewSkill("1234", "SomeName", ""))
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/skills", reader), &MockDataAccessor{})

	sc := SkillsController{BaseController: &base}
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to not specifying value for \"SkillType\" field in Skill POST request.")
	}
}

func TestPostSkillInvalidType(t *testing.T) {
	base := BaseController{}
	b, _ := json.Marshal(model.NewSkill("", "", "badtype"))
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/skills", reader), &MockDataAccessor{})

	sc := SkillsController{BaseController: &base}
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostNoSkill(t *testing.T) {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/skills/", nil), &MockDataAccessor{})
	sc := SkillsController{BaseController: &base}
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkillError(t *testing.T) {
	base := BaseController{}
	b, _ := json.Marshal(model.NewSkill("", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	base.Init(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/skills", reader), &MockErrorDataAccessor{})

	sc := SkillsController{BaseController: &base}
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}
