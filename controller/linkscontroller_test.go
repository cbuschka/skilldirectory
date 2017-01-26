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

func TestLinksControllerBase(t *testing.T) {
	base := BaseController{}
	lc := LinksController{BaseController: &base}

	if base != *lc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllLinks(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/links", nil)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllLinks_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/links", nil)
	lc := getLinksController(request, &data.MockErrorDataAccessor{})

	err := lc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/links/1234", nil)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetLink_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/links/1234", nil)
	lc := getLinksController(request, &data.MockErrorDataAccessor{})

	err := lc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/links/1234", nil)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteLink_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/links/1234", nil)
	lc := getLinksController(request, &data.MockErrorDataAccessor{})

	err := lc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteLink_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/links/", nil)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostLink(t *testing.T) {
	body := getReaderForNewLink("1234", "A Webpage", "http://webpage.com", "2345", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/links", body)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostLink_NoName(t *testing.T) {
	body := getReaderForNewLink("1234", "", "http://webpage.com", "2345", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/links", body)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "name")
	}
}

func TestPostLink_NoURL(t *testing.T) {
	body := getReaderForNewLink("1234", "A Webpage", "", "2345", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/links", body)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "url")
	}
}

func TestPostLink_NoSkillID(t *testing.T) {
	body := getReaderForNewLink("1234", "A Webpage", "http://webpage.com", "", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/links", body)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "skill_id")
	}
}

func TestPostLink_NoLinkType(t *testing.T) {
	body := getReaderForNewLink("1234", "A Webpage", "http://webpage.com", "2345", "")
	request := httptest.NewRequest(http.MethodPost, "/links", body)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "link_type")
	}
}

func TestPostLink_NoLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/links", nil)
	lc := getLinksController(request, &data.MockDataAccessor{})

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostLink_Error(t *testing.T) {
	body := getReaderForNewLink("1234", "A Webpage", "http://webpage.com", "2345", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/links", body)
	lc := getLinksController(request, &data.MockErrorDataAccessor{})

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func Test_validateLinkFields(t *testing.T) {
	lc := getLinksController(nil, &data.MockErrorDataAccessor{})
	link := model.Link{
		LinkType: model.WebpageLinkType,
		Name:     "Google",
		URL:      "http://www.google.com",
	}
	err := lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect empty Link.SkillID field.")
	}

	link = model.Link{
		SkillID: "1234",
		Name:    "Google",
		URL:     "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect empty Link.LinkType field.")
	}

	link = model.Link{
		SkillID:  "1234",
		LinkType: model.WebpageLinkType,
		URL:      "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect empty Link.Name field.")
	}

	link = model.Link{
		SkillID:  "1234",
		LinkType: model.WebpageLinkType,
		Name:     "Google",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect empty Link.URL field.")
	}

	link = model.Link{
		SkillID:  "1234",
		LinkType: model.WebpageLinkType,
		Name:     "Google",
		URL:      "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect invalid Link.SkillID.")
	}

	link = model.Link{
		SkillID:  "1234",
		LinkType: "MumboJumbo",
		Name:     "Google",
		URL:      "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect invalid Link.LinkType.")
	}
}

/*
getLinksController is a helper function for creating and initializing a new
BaseController with the given HTTP request and DataAccessor. Returns a new
LinksController created with that BaseController.
*/
func getLinksController(request *http.Request, dataAccessor data.DataAccess) LinksController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor, nil, logrus.New())
	return LinksController{BaseController: &base}
}

/*
getReaderForNewLink is a helper function for a new Link with the given id, name, url, skillID, and linkType.
This Link is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewLink(id, name, url, skillID, linkType string) *bytes.Reader {
	newLink := model.NewLink(id, name, url, skillID, linkType)
	b, _ := json.Marshal(newLink)
	return bytes.NewReader(b)
}
