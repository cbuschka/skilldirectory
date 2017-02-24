package controller

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"skilldirectory/data"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestSkillIconsControllerBase(t *testing.T) {
	base := BaseController{}
	sc := SkillIconsController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expecting Base() to return base pointer")
	}
}

func TestGetAllSkillIcons_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillicons", nil)
	sc := getSkillIconsController(request, nil, false)

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkillIcon_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/skillicons/1234", nil)
	sc := getSkillIconsController(request, data.MockFileSystem{}, false)

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkillIcon(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/1234", nil)
	sc := getSkillIconsController(request, &data.MockFileSystem{}, false)

	err := sc.Delete()
	if err != nil {
		t.Errorf("Expected no error, but got one: %s", err.Error())
	}
}

func TestDeleteSkillIcon_File_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/1234", nil)
	sc := getSkillIconsController(request, &data.MockErrorFileSystem{}, false)

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestDeleteSkillIcon_Gorm_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/1234", nil)
	sc := getSkillIconsController(request, &data.MockFileSystem{}, true)

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestDeleteSkillIcon_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/", nil)
	sc := getSkillIconsController(request, &data.MockFileSystem{}, false)

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key")
	}
}

func TestDeleteSkillIcon_BadID(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/skillicons/a", nil)
	sc := getSkillIconsController(request, &data.MockFileSystem{}, false)

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when a non uint key")
	}
}

func TestPostSkillIcon(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	file, _ := os.Open(path.Dir(wd) + "/resources/test.png")
	defer file.Close()

	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", file, http.MethodPost)
	sc := getSkillIconsController(req, &data.MockFileSystem{}, false)

	// Execute the test:
	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPutSkillIcon(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	file, _ := os.Open(path.Dir(wd) + "/resources/test.png")
	defer file.Close()

	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", file, http.MethodPut)
	sc := getSkillIconsController(req, &data.MockFileSystem{}, false)

	// Execute the test:
	err := sc.Put()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkillIconBadId(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	file, _ := os.Open(path.Dir(wd) + "/resources/test.png")
	defer file.Close()

	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("a", file, http.MethodPost)
	sc := getSkillIconsController(req, &data.MockFileSystem{}, false)

	// Execute the test:
	err := sc.Post()
	if err == nil {
		t.Errorf("Expecting error for non-integer id")
	}
}

func TestPostSkillIcon_InvalidImage(t *testing.T) {
	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", nil, http.MethodPost)
	sc := getSkillIconsController(req, &data.MockFileSystem{}, false)

	// Execute the test:
	err := sc.Post()
	if err == nil {
		t.Errorf("Post failed to detect invalid image")
	}
}

func TestPostSkillIcon_File_Error(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	file, _ := os.Open(path.Dir(wd) + "/api/resources/test.png")
	defer file.Close()

	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", file, http.MethodPost)
	sc := getSkillIconsController(req, &data.MockErrorFileSystem{}, false)

	// Execute the test:
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestPostSkillIcon_Gorm_Error(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	file, _ := os.Open(path.Dir(wd) + "/resources/test.png")
	defer file.Close()

	// Get multipart POST request using test image; and new controller w/ mockers
	req, _ := newSkillIconPostRequest("1234", file, http.MethodPost)
	sc := getSkillIconsController(req, &data.MockFileSystem{}, true)

	// Execute the test:
	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestSkillIconsOptions(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/skillicons", nil)
	sic := getSkillIconsController(request, nil, false)

	err := sic.Options()
	if err != nil {
		t.Errorf("OPTIONS requests should always return a 200 response.")
	}
	if sic.w.Header().Get("Access-Control-Allow-Methods") != "PUT, "+GetDefaultMethods() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Methods' contains" +
			" incorrect value")
	}
	if sic.w.Header().Get("Access-Control-Allow-Headers") != GetDefaultHeaders() {
		t.Errorf("OPTIONS response header 'Access-Control-Allow-Headers' contains" +
			" incorrect value")
	}
}

/*
getSkillIconsController is a helper function for creating and initializing a new
BaseController with the given HTTP request and err bool. Returns a new
SkillIconsController created with that BaseController.
*/
func getSkillIconsController(request *http.Request, fileSystem data.FileSystem, errSwitch bool) SkillIconsController {
	base := BaseController{}
	base.SetTest(errSwitch)
	base.Init(httptest.NewRecorder(), request, nil, fileSystem, logrus.New())
	return SkillIconsController{BaseController: &base}
}

func newSkillIconPostRequest(skillID string, icon *os.File, method string) (*http.Request, error) {
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

	req := httptest.NewRequest(method, "/api/skillicons", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
