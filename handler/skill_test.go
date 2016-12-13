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

func TestRemoveSkillNoID(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	b, _ := json.Marshal(model.NewSkill("", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	r := httptest.NewRequest(http.MethodDelete, "/skills", reader)
	err := removeSkill(r)
	if err == nil {
		t.Errorf("Expected error due to no ID in request.")
	}
}

func TestRemoveSkillErrorBadID(t *testing.T) {
	skillsConnector = data.NewAccessor(MockErrorDataAccessor{})
	b, _ := json.Marshal(model.NewSkill("", "", model.ScriptedSkillType))
	reader := bytes.NewReader(b)
	r := httptest.NewRequest(http.MethodDelete, "/skills/9000", reader)
	err := removeSkill(r)
	if err == nil {
		t.Errorf("Expected error due to non-existent skill ID in request.")
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
		t.Errorf("Expecting error for TestPerformGetError")
	}
}

func TestGetSkillsFiltered(t *testing.T) {
	skillsConnector = data.NewAccessor(NewMockInMemoryDataAccessor())

	newScriptedSkill := model.NewSkill("1234", "TestSkillName", model.ScriptedSkillType)
	skillsConnector.Save(newScriptedSkill.Id, newScriptedSkill)
	newCompiledSkill := model.NewSkill("2136", "TestSkillName", model.CompiledSkillType)
	skillsConnector.Save(newCompiledSkill.Id, newCompiledSkill)

	r := httptest.NewRequest(http.MethodGet, "/skills?skilltype=scripted", nil)
	w := httptest.NewRecorder()
	err := performGet(w, r)
	if err != nil {
		t.Errorf("Did not expect error when getting skills with filter")
	}

	correctResponseBody := "[{\"Id\":\"1234\",\"Name\":\"TestSkillName\",\"SkillType\":\"scripted\"}]"
	if w.Body.String() != correctResponseBody {
		t.Errorf("Failed to properly filter based on skilltype. "+
			"Expected Response body to be \n\t %s\n But got\n\t %s\\n",
			correctResponseBody, w.Body.String())
	}
}

func TestGetSkillsFilteredBadSkillType(t *testing.T) {
	skillsConnector = data.NewAccessor(MockDataAccessor{})
	r := httptest.NewRequest(http.MethodGet, "/skills?skilltype=badtype", nil)

	w := httptest.NewRecorder()
	err := performGet(w, r)
	if err == nil {
		t.Errorf("Expected error due to invalid skill type")
	}
}

func TestGetSkillsFilteredError(t *testing.T) {
	skillsConnector = data.NewAccessor(MockErrorDataAccessor{})
	r := httptest.NewRequest(http.MethodGet, "/skills?skilltype=scripted", nil)

	w := httptest.NewRecorder()
	err := performGet(w, r)
	if err == nil {
		t.Errorf("Expecting error for TestGetSkillsFilteredError")
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
