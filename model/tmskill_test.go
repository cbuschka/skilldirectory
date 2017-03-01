package model

import (
	"reflect"
	"testing"
)

func TestTMSkill_SetProficiency(t *testing.T) {
	tmSkill := NewTMSkillDefaults(1, 2, 3)

	tmSkill.SetProficiency(9000)
	if tmSkill.Proficiency != 5 {
		t.Error("method TMSkill.SetProficiency() failed to cap proficiency > 5")
	}
}

// Make sure the method works properly at the extremes
func TestTMSkill_GetProficiencyString(t *testing.T) {
	tmSkill := NewTMSkillDefaults(1, 2, 3)

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
	tms := NewTMSkillDefaults(0, 0, 0)
	if !reflect.DeepEqual(tms.GetType(), TMSkill{}) {
		t.Error("TMSkill getType not returning empty tmskill")
	}
}
