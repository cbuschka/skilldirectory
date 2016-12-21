package model

import (
	"reflect"
	"testing"
)

func TestNewSkill(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType).NewSkillDTO([]Link{NewLink("link", ".com", "ASkillID", BlogLinkType)})
	skillTwo := SkillDTO{
		Skill: Skill{
			ID:        "ASkillID",
			Name:      "ASkillName",
			SkillType: ScriptedSkillType},
		Links: []Link{NewLink("link", ".com", "ASkillID", BlogLinkType)},
	}
	// Verify that all of skillOne and skillTwo's fields are equal
	if !reflect.DeepEqual(skillOne, skillTwo) {
		t.Errorf("model/Skill\".NewSkill()\" produced incorrect Skill.")
	}
}

func TestSkillAddLink(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType).NewSkillDTO(nil)
	skillOne.AddLink(NewLink("Google", "http://www.google.com", skillOne.ID, WebpageLinkType))
	skillTwo := SkillDTO{
		Skill: Skill{ID: "ASkillID",
			Name:      "ASkillName",
			SkillType: ScriptedSkillType,
		},
		Links: []Link{NewLink("Google", "http://www.google.com", skillOne.ID, WebpageLinkType)},
	}

	// Verify that all of skillOne and skillTwo's fields are equal
	if !reflect.DeepEqual(skillOne, skillTwo) {
		t.Errorf("model/Skill\".AddLink()\" didn't work.")
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
