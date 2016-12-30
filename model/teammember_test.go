package model

import (
	"reflect"
	"testing"
)

func TestNewTeamMember(t *testing.T) {
	tmDTOOne := NewTeamMember("1234", "Yogi Bear", "Smarter Than Avg").NewTeamMemberDTO(
		[]TMSkill{NewTMSkillDefaults("2345", "3456", "1234")})
	tmDTOTwo := TeamMemberDTO{
		TeamMember: TeamMember{
			ID:    "1234",
			Name:  "Yogi Bear",
			Title: "Smarter Than Average Bear",
		},
		TMSkills: []TMSkill{NewTMSkillDefaults("2345", "3456", "1234")},
	}
	if reflect.DeepEqual(tmDTOOne, tmDTOTwo) {
		t.Error("model.teammember.\"NewTeamMemberDTO()\" produced incorrect TeamMemberDTO")
	}
}
