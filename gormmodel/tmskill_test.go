package gormmodel

import (
	"reflect"
	"testing"
)

func TestNewTMSkillDefaults(t *testing.T) {
	tmSkillDTOOne := NewTMSkillDefaults("TMSkillID", "SkillID",
		"TeamMemberID").NewTMSkillDTO("Coding", "Javatar")
	tmSkillDTOTwo := TMSkillDTO{
		TMSkill: TMSkill{
			ID:           "TMSkillID",
			SkillID:      "SkillID",
			TeamMemberID: "TeamMemberID",
			WishList:     false,
			Proficiency:  0,
		},
		SkillName:      "Coding",
		TeamMemberName: "Javatar",
	}
	// Verify that all of tmSkillOne and tmSkillTwo's fields are equal
	if !reflect.DeepEqual(tmSkillDTOOne, tmSkillDTOTwo) {
		t.Error("constructor NewTMSkillDefaults() produced incorrect TMSkill.")
	}
}

func TestNewTMSkillSetDefaults(t *testing.T) {
	tmSkillDTOOne := NewTMSkillSetDefaults("TMSkillID", "SkillID", "TeamMemberID",
		true, 3).NewTMSkillDTO("Coding", "Javatar")
	tmSkillDTOTwo := TMSkillDTO{
		TMSkill: TMSkill{
			ID:           "TMSkillId",
			SkillID:      "SkillID",
			TeamMemberID: "TeamMemberID",
			WishList:     true,
			Proficiency:  3,
		},
		SkillName:      "Coding",
		TeamMemberName: "Javatar",
	}
	// Verify that all of tmSkillOneDTO and tmSkillTwoDTO's fields are equal.
	if tmSkillDTOOne == tmSkillDTOTwo {
		t.Error("constructor NewTMSkillSetDefaults() produced incorrect TMSkill.")
	}

	tmSkill := NewTMSkillSetDefaults("TMSkillID", "SkillID", "TeamMemberID", true, 3)
	// Verify that the constructor clips proficiencies > 5 to 5
	tmSkill = NewTMSkillSetDefaults("TMSkillID", "SkillID", "TeamMemberID", true, 9000)
	if tmSkill.Proficiency != 5 {
		t.Error("constructor NewTMSkillSetDefaults() failed to cap proficiency > 5")
	}

	// Verify that the constructor clips proficiencies < 0 to 0
	tmSkill = NewTMSkillSetDefaults("TMSkillID", "SkillID", "TeamMemberID", true, -9000)
	if tmSkill.Proficiency != 0 {
		t.Error("constructor NewTMSkillSetDefaults() failed to cap proficiency < 0")
	}
}

func TestTMSkill_SetProficiency(t *testing.T) {
	tmSkill := NewTMSkillDefaults("TMSkillID", "SkillID", "TeamMemberID")

	tmSkill.SetProficiency(9000)
	if tmSkill.Proficiency != 5 {
		t.Error("method TMSkill.SetProficiency() failed to cap proficiency > 5")
	}

	tmSkill.SetProficiency(-9000)
	if tmSkill.Proficiency != 0 {
		t.Error("method TMSkill.SetProficiency() failed to cap proficiency < 0")
	}
}

// Make sure the method works properly at the extremes
func TestTMSkill_GetProficiencyString(t *testing.T) {
	tmSkill := NewTMSkillDefaults("TMSkillID", "SkillID", "TeamMemberID")

	tmSkill.SetProficiency(0)
	if tmSkill.GetProficiencyString() != "Not Applicable" {
		t.Error("method TMSkill.GetProficiency() failed to give correct " +
			"representation for Proficiency level 0.")
	}

	tmSkill.SetProficiency(1)
	if tmSkill.GetProficiencyString() != "Fundamentally Aware" {
		t.Error("method TMSkill.GetProficiency() failed to give correct " +
			"representation for Proficiency level 1.")
	}

	tmSkill.SetProficiency(2)
	if tmSkill.GetProficiencyString() != "Novice" {
		t.Error("method TMSkill.GetProficiency() failed to give correct " +
			"representation for Proficiency level 2.")
	}

	tmSkill.SetProficiency(3)
	if tmSkill.GetProficiencyString() != "Intermediate" {
		t.Error("method TMSkill.GetProficiency() failed to give correct " +
			"representation for Proficiency level 3.")
	}

	tmSkill.SetProficiency(4)
	if tmSkill.GetProficiencyString() != "Advanced" {
		t.Error("method TMSkill.GetProficiency() failed to give correct " +
			"representation for Proficiency level 4.")
	}

	tmSkill.SetProficiency(5)
	if tmSkill.GetProficiencyString() != "Expert" {
		t.Error("method TMSkill.GetProficiency() failed to give correct " +
			"representation for Proficiency level 5.")
	}

	tmSkill.Proficiency = 9000
	if tmSkill.GetProficiencyString() != "No String Representation Available" {
		t.Error("method TMSkill.GetProficiency() failed to give correct " +
			"representation for Proficiency level outside range 0-5.")
	}
}

func TestGetTMSkillType(t *testing.T) {
	tms := NewTMSkillDefaults("", "", "")
	if !reflect.DeepEqual(tms.GetType(), TMSkill{}) {
		t.Error("TMSkill getType not returning empty tmskill")
	}
}
