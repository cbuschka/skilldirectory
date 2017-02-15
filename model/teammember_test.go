package model

import (
	"reflect"
	"testing"
)

func TestNewTeamMember(t *testing.T) {
	tmOne := NewTeamMember("Yogi Bear", "Smarter Than Avg")
	tmTwo := TeamMember{
			Name:  "Yogi Bear",
			Title: "Smarter Than Average Bear",
	}
	if reflect.DeepEqual(tmOne, tmTwo) {
		t.Error("model.\"NewTeamMemberDTO()\" produced incorrect TeamMemberDTO")
	}
}

func TestGetTeamMemberType(t *testing.T) {
	tm := NewTeamMember("", "")
	if !reflect.DeepEqual(tm.GetType(), TeamMember{}) {
		t.Error("TeamMember getType not returning empty team member")
	}
}
