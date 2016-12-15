package model

import (
	"testing"
)

func TestNewSkill(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType)
	skillTwo := Skill{
		ID:        "ASkillID",
		Name:      "ASkillName",
		SkillType: ScriptedSkillType,
		Webpage:   Link{},
		Blogs:     []Link{},
		Tutorials: []Link{},
	}

	// Verify that all of skillOne and skillTwo's fields are equal
	if skillOne.ID != skillTwo.ID ||
		skillOne.Name != skillTwo.Name ||
		skillOne.SkillType != skillTwo.SkillType ||
		skillOne.Webpage != skillTwo.Webpage ||
		!compareLinkSlices(skillOne.Blogs, skillTwo.Blogs) ||
		!compareLinkSlices(skillOne.Tutorials, skillTwo.Tutorials) {
		t.Errorf("model/Skill\".NewSkill()\" produced incorrect Skill.")
	}
}

func TestNewSkillWithLinks(t *testing.T) {
	skillOne := NewSkillWithLinks("ASkillID", "ASkillName", ScriptedSkillType,
		Link{}, []Link{}, []Link{})
	skillTwo := Skill{
		ID:        "ASkillID",
		Name:      "ASkillName",
		SkillType: ScriptedSkillType,
		Webpage:   Link{},
		Blogs:     []Link{},
		Tutorials: []Link{},
	}

	// Verify that all of skillOne and skillTwo's fields are equal
	if skillOne.ID != skillTwo.ID ||
		skillOne.Name != skillTwo.Name ||
		skillOne.SkillType != skillTwo.SkillType ||
		skillOne.Webpage != skillTwo.Webpage ||
		!compareLinkSlices(skillOne.Blogs, skillTwo.Blogs) ||
		!compareLinkSlices(skillOne.Tutorials, skillTwo.Tutorials) {
		t.Errorf("model/Skill\".NewSkillWithLinks()\" produced incorrect Skill.")
	}
}

func TestSkillAddLink(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType)
	skillOne.AddLink(NewLink("Google", "http://www.google.com", skillOne.ID, WebpageLinkType))
	skillOne.AddLink(NewLink("Google Blog", "https://www.blog.google/", skillOne.ID, BlogLinkType))
	skillOne.AddLink(NewLink("Google Apps Script", "https://developers.google.com/apps-script/articles",
		skillOne.ID, TutorialLinkType))
	skillTwo := Skill{
		ID:        "ASkillID",
		Name:      "ASkillName",
		SkillType: ScriptedSkillType,
		Webpage:   NewLink("Google", "http://www.google.com", skillOne.ID, WebpageLinkType),
		Blogs:     []Link{NewLink("Google Blog", "https://www.blog.google/", skillOne.ID, BlogLinkType)},
		Tutorials: []Link{NewLink("Google Apps Script", "https://developers.google.com/apps-script/articles",
			skillOne.ID, TutorialLinkType)},
	}

	// Verify that all of skillOne and skillTwo's fields are equal
	if skillOne.ID != skillTwo.ID ||
		skillOne.Name != skillTwo.Name ||
		skillOne.SkillType != skillTwo.SkillType ||
		skillOne.Webpage != skillTwo.Webpage ||
		!compareLinkSlices(skillOne.Blogs, skillTwo.Blogs) ||
		!compareLinkSlices(skillOne.Tutorials, skillTwo.Tutorials) {
		t.Errorf("model/Skill\".AddLink()\" didn't work.")
	}
}

func TestSkillAddLinkBadLinkType(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType)

	invalidLinkType := "SillyLinkType"
	err := skillOne.AddLink(NewLink("Google", "http://www.google.com", skillOne.ID, invalidLinkType))
	if err == nil {
		t.Errorf("\"skill.AddLink()\" function failed to detect invalid LinkType: \"%s\"", invalidLinkType)
	}

	err = skillOne.AddLink(NewLink("Google", "http://www.google.com", skillOne.ID, ""))
	if err == nil {
		t.Errorf("\"skill.AddLink()\" function failed to detect empty (\"\") LinkType")
	}
}

func TestInvalidSkillType(t *testing.T) {
	if IsValidSkillType("InvalidSkillType") {
		t.Errorf("func IsValidSkillType() failed to detect invalid SkillType.")
	}
}

func TestValidSkillType(t *testing.T) {
	if !IsValidSkillType(ScriptedSkillType) {
		t.Errorf("func IsValidSkillType() flagged valid SkillType as invalid")
	}
}

func compareLinkSlices(a, b []Link) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
