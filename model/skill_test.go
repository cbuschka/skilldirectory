package model

import "testing"

func TestNewSkill(t *testing.T) {
	skillOne := NewSkill("ASkillID", "ASkillName", ScriptedSkillType).NewSkillWithLinks(Link{}, nil, nil)
	skillTwo := SkillDTO{
		Skill: Skill{
			ID:        "ASkillID",
			Name:      "ASkillName",
			SkillType: ScriptedSkillType},
		Webpage:   Link{},
		Blogs:     nil,
		Tutorials: nil,
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
