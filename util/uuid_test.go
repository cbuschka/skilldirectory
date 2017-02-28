package util

import (
	"net/url"
	"testing"
	"time"
)

func TestNewId(t *testing.T) {
	id1 := NewID()
	time.Sleep(1)
	id2 := NewID()
	if id1 == id2 {
		t.Error("Random id failed")
	}

}

func TestIdLength(t *testing.T) {
	id := NewID()
	if len(id) != 36 {
		t.Errorf("ID length incorrect")
	}
}

func TestCheckForID(t *testing.T) {
	url, _ := url.Parse("https://test.com/path")
	path := CheckForID(url)
	if path != "path" {
		t.Errorf("Expecting 'path' got: %s", path)
	}
}

func TestPathToID(t *testing.T) {
	url, _ := url.Parse("https://test.com/1")
	id, err := PathToID(url)
	if err != nil {
		t.Errorf("Error from PathToID: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected '1' got %d", id)
	}
}

func TestStringToID(t *testing.T) {
	path := "1"
	id, err := StringToID(path)
	if err != nil {
		t.Errorf("Error from PathToID: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected '1' got %d", id)
	}
}

func TestStringToID0(t *testing.T) {
	path := "-1"
	_, err := StringToID(path)
	if err == nil {
		t.Errorf("Expected Error from PathToID with a 0")
	}

}

func TestStringToIDNilID(t *testing.T) {
	path := ""
	_, err := StringToID(path)
	if err == nil {
		t.Errorf("Expected Error from PathToID with nil string")
	}

}
