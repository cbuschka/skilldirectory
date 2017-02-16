package util

import (
	"net/url"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestGetProperty(t *testing.T) {
	os.Setenv("testkey", "testvalue")
	value := GetProperty("testkey")
	if value != "testvalue" {
		t.Error("GetProperty not getting values")
	}
}

func TestCheckForId(t *testing.T) {
	url := url.URL{}
	url.Path = "/api/skills/"

	if CheckForID(&url) != "" {
		t.Errorf("No id failed")
	}

}

func TestCheckForIdError(t *testing.T) {
	url := url.URL{}
	url.Path = "/api/skils/"

	if CheckForID(&url) == "" {
		t.Errorf("No id failed")
	}
}

func TestURLIdParse(t *testing.T) {
	url := url.URL{}
	url.Path = "/api/skills/abc"

	if CheckForID(&url) != "abc" {
		t.Errorf("Id match failed")
	}
}

func TestRootDir(t *testing.T) {
	path := "skills/id"
	rootDir := getRootDir(path)
	if rootDir != "skills/" {
		t.Errorf("Rootdir parse failed")
	}
}

func TestRootDirNoId(t *testing.T) {
	path := "skills"
	rootDir := getRootDir(path)
	if rootDir != "skills/" {
		t.Errorf("Rootdir parse failed")
	}
}

func TestRootDirSlash(t *testing.T) {
	path := "/api/skills"
	rootDir := getRootDir(path)
	if rootDir != "skills/" {
		t.Errorf("Rootdir parse failed. Root = %s", rootDir)
	}
}

func TestValidateIcon(t *testing.T) {
	// Open test PNG image file
	wd, _ := os.Getwd()
	icon, _ := os.Open(path.Dir(wd) + "/resources/test.png")
	defer icon.Close()

	_, err := ValidateIcon(icon)
	if err != nil {
		t.Errorf("Flagged valid icon as invalid: %s", err)
	}
}

func TestValidateIconError(t *testing.T) {
	wd, _ := os.Getwd()
	icon, _ := os.Open(path.Dir(wd) + "/resources/est.png")
	defer icon.Close()

	_, err := ValidateIcon(icon)
	if err == nil {
		t.Errorf("Expected error with a nil Reader")
	}
}

func TestSanitizeInput(t *testing.T) {
	testString := "'test'"
	expectedString := "''test''"
	returnString := SanitizeInput(testString)
	if !reflect.DeepEqual(returnString, expectedString) {
		t.Errorf("Expected %s got %s", expectedString, returnString)
	}
}
