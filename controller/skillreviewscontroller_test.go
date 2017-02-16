package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"skilldirectory/data"
	"skilldirectory/model"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestSkillReviewsController_Base(t *testing.T) {
	base := BaseController{}
	sc := SkillReviewsController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllSkillReviews(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillreviews", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllSkillReviews_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillreviews", nil)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkillReview(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillreviews/1234", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSkillReview_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillreviews/1234", nil)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkillReview(t *testing.T) {
	body := getReaderForDeleteSkillReview("1234", "2345")
	request := httptest.NewRequest(http.MethodDelete, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Delete()
	if err != nil {
		t.Errorf("Did not expect error: %s", err.Error())
	}
}

func TestDeleteSkillReview_Error(t *testing.T) {
	body := getReaderForDeleteSkillReview("1234", "2345")
	request := httptest.NewRequest(http.MethodDelete, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkillReview(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah",
		"1234", true)
	request := httptest.NewRequest(http.MethodPost, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkillReview_NoSkillID(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "", "3456", "blah", "1234",
		true)
	request := httptest.NewRequest(http.MethodPost, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in SkillReview POST"+
			" request.", "skill_id")
	}
}

func TestPostSkillReview_NoTeamMemberID(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "", "blah", "1234",
		true)
	request := httptest.NewRequest(http.MethodPost, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in SkillReview POST"+
			" request.", "team_member_id")
	}
}

func TestPostSkillReview_NoBody(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "", "1234",
		true)
	request := httptest.NewRequest(http.MethodPost, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error dur to empty %q field in SkillReview POST"+
			" request.", "body")
	}
}

func TestPostSkillReview_NoSkillReview(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/skillreviews", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkillReview_Error(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah",
		"1234", true)
	request := httptest.NewRequest(http.MethodPost, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPutSkillReview(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah", "1234",
		true)
	request := httptest.NewRequest(http.MethodPut, "/api/skillreviews/1234", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Put()
	if err != nil {
		t.Errorf("Expected error due to empty %q field in SkillReview POST"+
			" request", "body")
	}
}

func TestPutSkillReviewNoId(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah", "1234",
		true)
	request := httptest.NewRequest(http.MethodPut, "/api/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Put()
	if err == nil {
		t.Errorf("Expected error due to no id in request path")
	}
}

func TestPutSkillReviewError(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah", "1234",
		true)
	request := httptest.NewRequest(http.MethodPut, "/api/skillreviews/1234", body)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Put()
	if err == nil {
		t.Errorf("Expected error due to no backend fail")
	}
}

func TestConvertToReviewsStruct(t *testing.T) {
	preReview := []interface{}{model.SkillReview{SkillID: "1234"}, model.SkillReview{SkillID: "5678"}}
	pretReviewStruct := []model.SkillReview{model.SkillReview{SkillID: "1234"}, model.SkillReview{SkillID: "5678"}}
	skillReview, err := convertToReviewsStruct(preReview)
	if err != nil {
		t.Errorf("Convert to struct failed: %v", err)
	}

	if !reflect.DeepEqual(pretReviewStruct, skillReview) {
		t.Error("Deep equal failed for convert review to struct")
	}
}

func TestConvertReviewsToDTOs(t *testing.T) {
	preReview := []model.SkillReview{
		model.SkillReview{SkillID: "1234"},
		model.SkillReview{SkillID: "5678"},
	}
	preReviewDTO := []model.SkillReviewDTO{
		preReview[0].NewSkillReviewDTO("", ""),
		preReview[1].NewSkillReviewDTO("", ""),
	}
	reviewsController := getSkillReviewsController(nil, data.MockDataAccessor{})
	reviewsDTO := reviewsController.convertReviewsToDTOs(preReview)
	if !reflect.DeepEqual(preReviewDTO, reviewsDTO) {
		t.Error("Expecting a match of tmskills -> tmskillsDTO")
	}
}

func TestConvertReviewsToDTOsError(t *testing.T) {
	preReview := []model.SkillReview{
		model.SkillReview{SkillID: "1234"},
		model.SkillReview{SkillID: "5678"},
	}
	preReviewDTO := []model.SkillReviewDTO{}
	reviewsController := getSkillReviewsController(nil, data.MockErrorDataAccessor{})
	reviewsDTO := reviewsController.convertReviewsToDTOs(preReview)
	if !reflect.DeepEqual(preReviewDTO, reviewsDTO) {
		t.Error("Expecting an array of 0 SkillReviews")
	}
}

/*
getSkillReviewsController is a helper function for creating and initializing a new
BaseController with the given HTTP request and DataAccessor. Returns a new
SkillReviewsController created with that BaseController.
*/
func getSkillReviewsController(request *http.Request,
	dataAccessor data.DataAccess) SkillReviewsController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor, nil, logrus.New())
	return SkillReviewsController{BaseController: &base}
}

/*
getReaderForNewSkillReview is a helper function for a new SkillReview with the given
id, skillID, teamMemberID, body, date, and positive flag. This SkillReview is then
marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewSkillReview(id, skillID, teamMemberID, body, timestamp string,
	positive bool) *bytes.Reader {
	newSkillReview := model.NewSkillReview(id, skillID, teamMemberID,
		body, timestamp, positive)
	b, _ := json.Marshal(newSkillReview)
	return bytes.NewReader(b)
}

func getReaderForDeleteSkillReview(id string, skillID string) *bytes.Reader {
	newSkillReview := model.SkillReview{
		ID:      id,
		SkillID: skillID,
	}
	b, _ := json.Marshal(newSkillReview)
	return bytes.NewReader(b)
}
