package model

import (
	"testing"
)

func TestNewSkill(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType)
	skillTwo := Skill{"ASkillID", "ASkillName", ScriptedSkillType}
	if skillOne != skillTwo {
		t.Errorf("model/Skill\".NewSkill()\" produced incorrect Skill.")
	}
}

func TestInValidSkillType(t *testing.T) {
	if IsValidSkillType("InvalidSkillType") {
		t.Errorf("func IsValidSkillType() failed to detect invalid SkillType.")
	}
}

func TestValidSkillType(t *testing.T) {
	if !IsValidSkillType(ScriptedSkillType) {
		t.Errorf("func IsValidSkillType() flagged valid SkillType as invalid")
	}
}