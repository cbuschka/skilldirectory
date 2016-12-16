package model

import "testing"

func TestNewLink(t *testing.T) {
	linkOne := NewLink("Google", "http://www.google.com", "1234", WebpageLinkType)
	linkTwo := Link{
		Name:     "Google",
		URL:      "http://www.google.com",
		SkillID:  "1234",
		LinkType: WebpageLinkType,
	}

	if linkOne != linkTwo {
		t.Errorf("model/link\".NewLink()\" produced incorrect Link.")
	}
}

func TestInvalidLinkType(t *testing.T) {
	if IsValidLinkType("InvalidLinkType") {
		t.Errorf("func IsValidLinkType() failed to detect invalid LinkType")
	}
}

func TestValidLinkType(t *testing.T) {
	if !IsValidLinkType(WebpageLinkType) {
		t.Errorf("func IsValidLinkType() flagged calid LinkType as invalid.")
	}
}