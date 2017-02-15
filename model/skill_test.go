package model

import (
	"reflect"
	"testing"
)

func TestNewSkill(t *testing.T) {
	skillOne := NewSkill("ASkillName", ScriptedSkillType, "example.com");
	skillTwo := Skill{
		Name: "ASkillName",
		SkillType: ScriptedSkillType,
		IconURL: "example.com",
	}
	// Verify that all of skillOne and skillTwo's fields are equal
	if !reflect.DeepEqual(skillOne, skillTwo) {
		t.Error("\"model.NewSkill()\" produced incorrect Skill.")
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
