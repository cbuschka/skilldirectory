package util

import (
	"net/url"
	"os"
	"path"
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
	url.Path = "/skills/"

	if CheckForID(&url) != "" {
		t.Errorf("No id failed")
	}

}

func TestCheckForIdError(t *testing.T) {
	url := url.URL{}
	url.Path = "/skils/"

	if CheckForID(&url) == "" {
		t.Errorf("No id failed")
	}
}

func TestURLIdParse(t *testing.T) {
	url := url.URL{}
	url.Path = "/skills/abc"

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
	path := "/skills"
	rootDir := getRootDir(path)
	if rootDir != "skills/" {
		t.Errorf("Rootdir parse failed")
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
