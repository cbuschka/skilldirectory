package controller

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path"
	"skilldirectory/data"
	"testing"

	"github.com/Sirupsen/logrus"

	"os"
)

func TestSkillIconsControllerBase(t *testing.T) {
	base := BaseController{}
	sc := SkillsController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expecting Base() to return base pointer")
	}
}

func TestGetAllSkillIcons(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillicons", nil)
	sc := getSkillIconsController(request, &data.MockDataAccessor{},
		&data.MockFileSystem{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllSkillIcons_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillicons", nil)
	sc := getSkillIconsController(request, &data.MockErrorDataAccessor{},
		&data.MockErrorFileSystem{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkillIcon(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillicon/1234", nil)
	sc := getSkillIconsController(request, &data.MockDataAccessor{},
		&data.MockFileSystem{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSkillIcon_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillicons/1234", nil)
	sc := getSkillIconsController(request, &data.MockErrorDataAccessor{},
		&data.MockErrorFileSystem{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkillIcon(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/1234", nil)
	sc := getSkillIconsController(request, &data.MockDataAccessor{},
		&data.MockFileSystem{})

	err := sc.Delete()
	if err != nil {
		t.Errorf("Expected no error, but got one: %s", err.Error())
	}
}

func TestDeleteSkillIcon_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/1234", nil)
	sc := getSkillIconsController(request, &data.MockErrorDataAccessor{},
		&data.MockErrorFileSystem{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkillIcon_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/", nil)
	sc := getSkillIconsController(request, &data.MockDataAccessor{},
		&data.MockFileSystem{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostSkillIcon(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	file, _ := os.Open(path.Dir(wd) + "/resources/test.png")
	defer file.Close()

	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", file)
	sc := getSkillIconsController(req, &data.MockDataAccessor{}, &data.MockFileSystem{})

	// Execute the test:
	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkillIcon_InvalidImage(t *testing.T) {
	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", nil)
	sc := getSkillIconsController(req, &data.MockDataAccessor{}, &data.MockFileSystem{})

	// Execute the test:
	err := sc.Post()
	if err == nil {
		t.Errorf("Post failed to detect invalid image: %s", err.Error())
	}
}

func TestPostSkillIcon_Error(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	file, _ := os.Open(path.Dir(wd) + "/api/resources/test.png")
	defer file.Close()

	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", file)
	sc := getSkillIconsController(req, &data.MockErrorDataAccessor{},
		&data.MockErrorFileSystem{})

	// Execute the test:
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

/*
getSkillIconsController is a helper function for creating and initializing a new
BaseController with the given HTTP request and DataAccessor. Returns a new
SkillsController created with that BaseController.
*/
func getSkillIconsController(request *http.Request, dataAccessor data.DataAccess,
	fileSystem data.FileSystem) SkillIconsController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor, fileSystem, logrus.New())
	return SkillIconsController{BaseController: &base}
}

func newSkillIconPostRequest(skillID string, icon *os.File) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if icon != nil {
		part, err := writer.CreateFormFile("icon", "test.png")
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, icon)
	}

	writer.WriteField("skill_id", skillID)
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/api/skillicons", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
