package model

import (
	"reflect"
	"testing"
)

func TestLink_NewLink(t *testing.T) {
	linkOne := NewLink(1234, 2345, "Google", "http://www.google.com", WebpageLinkType)
	linkTwo := Link{
		Name:     "Google",
		URL:      "http://www.google.com",
		SkillID:  2345,
		LinkType: WebpageLinkType,
	}
	linkTwo.ID = 1234

	if linkOne != linkTwo {
		t.Errorf("\"model.NewLink()\" produced incorrect Link.")
	}
}

func TestIsValidLinkType1(t *testing.T) {
	if IsValidLinkType("InvalidLinkType") {
		t.Errorf("func IsValidLinkType() failed to detect invalid LinkType")
	}
}

func TestIsValidLinkType2(t *testing.T) {
	if !IsValidLinkType(WebpageLinkType) {
		t.Errorf("func IsValidLinkType() flagged valid LinkType as invalid.")
	}
}

func TestGetLinkType(t *testing.T) {
	l := NewLink(0, 0, "", "", "")
	if !reflect.DeepEqual(l.GetType(), Link{}) {
		t.Error("Link getType not returning empty link")
	}
}
