package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/data"
	"skilldirectory/model"
	"testing"
)

func TestSkillReviewsController_Base(t *testing.T) {
	base := BaseController{}
	sc := SkillReviewsController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllSkillReviews(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skillreviews", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllSkillReviews_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skillreviews", nil)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkillReview(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skillreviews/1234", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSkillReview_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skillreviews/1234", nil)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkillReview(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/skillreviews/1234", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkillReview_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/skillreviews/1234", nil)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkillReview_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/skillreviews/", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostSkillReview(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah",
		"12/28/2016", true)
	request := httptest.NewRequest(http.MethodPost, "/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkillReview_NoSkillID(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "", "3456", "blah", "12/28/2016",
		true)
	request := httptest.NewRequest(http.MethodPost, "/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in SkillReview POST"+
			" request.", "skill_id")
	}
}

func TestPostSkillReview_NoTeamMemberID(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "", "blah", "12/28/2016",
		true)
	request := httptest.NewRequest(http.MethodPost, "/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in SkillReview POST"+
			" request.", "team_member_id")
	}
}

func TestPostSkillReview_NoBody(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "", "12/28/2016",
		true)
	request := httptest.NewRequest(http.MethodPost, "/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error dur to empty %q field in SkillReview POST"+
			" request.", "body")
	}
}

func TestPostSkillReview_NoDate(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah", "", true)
	request := httptest.NewRequest(http.MethodPost, "/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error dur to empty %q field in SkillReview POST"+
			" request.", "date")
	}
}

func TestPostSkillReview_NoSkillReview(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/skillreviews", nil)
	sc := getSkillReviewsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkillReview_Error(t *testing.T) {
	body := getReaderForNewSkillReview("1234", "2345", "3456", "blah",
		"12/28/201", true)
	request := httptest.NewRequest(http.MethodPost, "/skillreviews", body)
	sc := getSkillReviewsController(request, &data.MockErrorDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
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
	base.Init(httptest.NewRecorder(), request, dataAccessor)
	return SkillReviewsController{BaseController: &base}
}

/*
getReaderForNewSkillReview is a helper function for a new SkillReview with the given
id, skillID, teamMemberID, body, date, and positive flag. This SkillReview is then
marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewSkillReview(id, skillID, teamMemberID, body, date string,
	positive bool) *bytes.Reader {
	newSkillReview := model.NewSkillReview(id, skillID, teamMemberID,
		body, date, positive)
	b, _ := json.Marshal(newSkillReview)
	return bytes.NewReader(b)
}
