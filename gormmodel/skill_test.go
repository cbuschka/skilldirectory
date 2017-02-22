package gormmodel

import (
	"reflect"
	"testing"
)

func TestNewSkill(t *testing.T) {
	skillOne := NewSkill(1, "ASkillName", ScriptedSkillType)
	skillTwo := Skill{
		Name:      "ASkillName",
		SkillType: ScriptedSkillType}
	skillTwo.ID = 1
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
	s := NewSkill(1, "", "")
	if !reflect.DeepEqual(s.GetType(), Skill{}) {
		t.Error("Skill getType not returning empty skill")
	}
}
