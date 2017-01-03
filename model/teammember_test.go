package model

import (
	"reflect"
	"testing"
)

func TestNewTeamMember(t *testing.T) {
	tmDTOOne := NewTeamMember("1234", "Yogi Bear", "Smarter Than Avg").
		NewTeamMemberDTO([]TMSkillDTO{
			NewTMSkillDefaults("2345", "3456", "1234").
				NewTMSkillDTO("whatever", "Yogi Bear")})
	tmDTOTwo := TeamMemberDTO{
		TeamMember: TeamMember{
			ID:    "1234",
			Name:  "Yogi Bear",
			Title: "Smarter Than Average Bear",
		},
		TMSkillDTOs: []TMSkillDTO{
			NewTMSkillDefaults("2345", "3456", "1234").
				NewTMSkillDTO("whatever", "Yogi Bear")},
	}
	if reflect.DeepEqual(tmDTOOne, tmDTOTwo) {
		t.Error("model.\"NewTeamMemberDTO()\" produced incorrect TeamMemberDTO")
	}
}
