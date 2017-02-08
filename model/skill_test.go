package model

import (
	"reflect"
	"testing"
)

func TestNewSkill(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType).NewSkillDTO(
		[]Link{NewLink("1234", "link", ".com", "ASkillID", BlogLinkType)},
		NewSkillIcon("1234", "http://abc.xyz/123.png"))
	skillTwo := SkillDTO{
		Skill: Skill{
			ID:        "ASkillID",
			Name:      "ASkillName",
			SkillType: ScriptedSkillType},
		Links: []Link{NewLink("1234", "link", ".com", "ASkillID", BlogLinkType)},
		Icon:  NewSkillIcon("1234", "http://abc.xyz/123.png"),
	}
	// Verify that all of skillOne and skillTwo's fields are equal
	if !reflect.DeepEqual(skillOne, skillTwo) {
		t.Error("\"model.NewSkill()\" produced incorrect Skill.")
	}
}

func TestSkillAddLink(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType).NewSkillDTO(
		nil, NewSkillIcon("1234", "http://abc.xyz/123.png"))
	skillOne.AddLink(NewLink("1234", "Google", "http://www.google.com",
		skillOne.ID, WebpageLinkType))
	skillTwo := SkillDTO{
		Skill: Skill{ID: "ASkillID",
			Name:      "ASkillName",
			SkillType: ScriptedSkillType,
		},
		Links: []Link{NewLink("1234", "Google", "http://www.google.com", skillOne.ID,
			WebpageLinkType)},
		Icon: NewSkillIcon("1234", "http://abc.xyz/123.png"),
	}

	// Verify that all of skillOne and skillTwo's fields are equal
	if !reflect.DeepEqual(skillOne, skillTwo) {
		t.Error("model.Skill.AddLink() didn't work.")
	}
}

func TestInvalidSkillType(t *testing.T) {
	if IsValidSkillType("InvalidSkillType") {
		t.Error("func IsValidSkillType() failed to detect invalid SkillType.")
	}
}

func TestValidSkillType(t *testing.T) {
	if !IsValidSkillType(ScriptedSkillType) {
		t.Error("func IsValidSkillType() flagged valid SkillType as invalid")
	}
}

func TestGetSkillType(t *testing.T) {
	s := NewSkill("", "", "")
	if !reflect.DeepEqual(s.GetType(), Skill{}) {
		t.Error("Skill getType not returning empty skill")
	}
}
