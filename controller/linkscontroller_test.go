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

func TestLinksControllerBase(t *testing.T) {
	base := BaseController{}
	lc := LinksController{BaseController: &base}

	if base != *lc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllLinks(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/links", nil)
	lc := getLinksController(request, false)

	err := lc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllLinksLinkTypeFilter(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/links?linktype=blog", nil)
	lc := getLinksController(request, false)

	err := lc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllLinksLinkTypeFilterError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/links?linktype=blog", nil)
	lc := getLinksController(request, true)

	err := lc.Get()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestGetAllLinks_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/links", nil)
	lc := getLinksController(request, true)

	err := lc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/links/1234", nil)
	lc := getLinksController(request, false)

	err := lc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetLink_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/links/1234", nil)
	lc := getLinksController(request, true)

	err := lc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/links/1234", nil)
	lc := getLinksController(request, false)

	err := lc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteLink_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/links/1234", nil)
	lc := getLinksController(request, true)

	err := lc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteLink_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/links/", nil)
	lc := getLinksController(request, false)

	err := lc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostLink(t *testing.T) {
	body := getReaderForNewLink(1234, 2345, "A Webpage", "http://webpage.com", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/api/links", body)
	lc := getLinksController(request, false)

	err := lc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostLink_NoName(t *testing.T) {
	body := getReaderForNewLink(1234, 2345, "", "http://webpage.com", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/api/links", body)
	lc := getLinksController(request, false)

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "name")
	}
}

func TestPostLink_NoURL(t *testing.T) {
	body := getReaderForNewLink(1234, 2345, "A Webpage", "", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/api/links", body)
	lc := getLinksController(request, false)

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "url")
	}
}

func TestPostLink_NoSkillID(t *testing.T) {
	body := getReaderForNewLink(1234, 0, "A Webpage", "http://webpage.com", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/api/links", body)
	lc := getLinksController(request, false)

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "skill_id")
	}
}

func TestPostLink_NoLinkType(t *testing.T) {
	body := getReaderForNewLink(1234, 2345, "A Webpage", "http://webpage.com", "")
	request := httptest.NewRequest(http.MethodPost, "/api/links", body)
	lc := getLinksController(request, false)

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Link POST request.", "link_type")
	}
}

func TestPostLink_NoLink(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/links", nil)
	lc := getLinksController(request, false)

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostLink_Error(t *testing.T) {
	body := getReaderForNewLink(1234, 2345, "A Webpage", "http://webpage.com", model.WebpageLinkType)
	request := httptest.NewRequest(http.MethodPost, "/api/links", body)
	lc := getLinksController(request, true)

	err := lc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func Test_validateLinkFields(t *testing.T) {
	lc := getLinksController(nil, false)
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
		SkillID: 1234,
		Name:    "Google",
		URL:     "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect empty Link.LinkType field.")
	}

	link = model.Link{
		SkillID:  1234,
		LinkType: model.WebpageLinkType,
		URL:      "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect empty Link.Name field.")
	}

	link = model.Link{
		SkillID:  1234,
		LinkType: model.WebpageLinkType,
		Name:     "Google",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect empty Link.URL field.")
	}

	link = model.Link{
		SkillID:  0,
		LinkType: model.WebpageLinkType,
		Name:     "Google",
		URL:      "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect invalid Link.SkillID.")
	}

	link = model.Link{
		SkillID:  1234,
		LinkType: "MumboJumbo",
		Name:     "Google",
		URL:      "http://www.google.com",
	}
	err = lc.validateLinkFields(&link)
	if err == nil {
		t.Errorf("validateLinkFields() failed to detect invalid Link.LinkType.")
	}
}

func TestLinksOptions(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/links", nil)
	lc := getLinksController(request, false)

	err := lc.Options()
	if err != nil {
		t.Errorf("OPTIONS requests should always return a 200 response.")
	}
	if lc.w.Header().Get("Access-Control-Allow-Methods") != GetDefaultMethods() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Methods' contains" +
			" incorrect value")
	}
	if lc.w.Header().Get("Access-Control-Allow-Headers") != GetDefaultHeaders() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Headers' contains" +
			" incorrect value")
	}
}

/*
getLinksController is a helper function for creating and initializing a new
BaseController with the given HTTP request and DataAccessor. Returns a new
LinksController created with that BaseController.
*/
func getLinksController(request *http.Request, errSwitch bool) LinksController {
	base := BaseController{}
	base.SetTest(errSwitch)
	base.InitWithGorm(httptest.NewRecorder(), request, nil, logrus.New(), nil)
	return LinksController{BaseController: &base}
}

/*
getReaderForNewLink is a helper function for a new Link with the given id, name, url, skillID, and linkType.
This Link is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewLink(id, skillID uint, name, url, linkType string) *bytes.Reader {
	newLink := model.NewLink(id, skillID, name, url, linkType)
	b, _ := json.Marshal(newLink)

	return bytes.NewReader(b)
}
