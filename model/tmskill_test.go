package model

import (
	"testing"
	"fmt"
)

func TestNewTMSkillDefaults(t *testing.T) {
	tmSkillOne := NewTMSkillDefaults("TMSkillID", "SkillID", "TeamMemberID")
	tmSkillTwo := TMSkill{
		ID:           "TMSkillID",
		SkillID:      "SkillID",
		TeamMemberID: "TeamMemberID",
		WishList:     false,
		Proficiency:  0,
	}
	//Verify that all of tmSkillOne and tmSkillTwo's fields are equal
	if tmSkillOne != tmSkillTwo {
		t.Errorf("constructor newTMSkillDefaults() produced incorrect TMSkill.")
	}
}

func TestNewTMSkillSetDefaults(t *testing.T) {
	tmSkillOne := NewTMSkillSetDefaults("TMSkillID", "SkillID", "TeamMemberID", true, 3)
	tmSkillTwo := TMSkill{
		ID:           "TMSkillID",
		SkillID:      "SkillID",
		TeamMemberID: "TeamMemberID",
		WishList:     true,
		Proficiency:  3,
	}
	// Verify that all of tmSkillOne and tmSkillTwo's fields are equal.
	if tmSkillOne != tmSkillTwo {
		t.Errorf("constructor newTMSkillSetDefaults() produced incorrect TMSkill.")
	}

	// Verify that the constructor clips proficiencies > 5 to 5
	tmSkillOne = NewTMSkillSetDefaults("TMSkillID", "SkillID", "TeamMemberID", true, 9000)
	if tmSkillOne.Proficiency != 5 {
		t.Error("constructor newTMSkillSetDefaults() failed to cap proficiency > 5")
	}

	// Verify that the constructor clips proficiencies < 0 to 0
	tmSkillOne = NewTMSkillSetDefaults("TMSkillID", "SkillID", "TeamMemberID", true, -9000)
	if tmSkillOne.Proficiency != 0 {
		t.Error("constructor newTMSkillSetDefaults() failed to cap proficiency < 0")
	}
}

func TestTMSkill_SetProficiency(t *testing.T) {
	tmSkill := NewTMSkillDefaults("TMSkillID", "SkillID", "TeamMemberID")

	tmSkill.SetProficiency(9000)
	if tmSkill.Proficiency != 5 {
		t.Error("method TMSkill.setProficiency() failed to cap proficiency > 5")
	}

	tmSkill.SetProficiency(-9000)
	if tmSkill.Proficiency != 0 {
		t.Error("method TMSkill.setProficiency() failed to cap proficiency < 0")
	}
}
