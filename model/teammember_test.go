package model

import (
	"reflect"
	"testing"
)

func TestNewTeamMember(t *testing.T) {

	teamMemberOne := NewTeamMember(1234, "Yogi Bear", "Smarter Than Avg")
	teamMemberTwo := TeamMember{
		Name:  "Yogi Bear",
		Title: "Smarter Than Avg",
	}
	teamMemberTwo.ID = 1234
	if !reflect.DeepEqual(teamMemberOne, teamMemberTwo) {
		t.Errorf("New TeamMember Error.\nOne: %v \nTwo: %v", teamMemberOne, teamMemberTwo)
	}
}

func TestGetTeamMemberType(t *testing.T) {
	tm := NewTeamMember(0, "", "")
	if !reflect.DeepEqual(tm.GetType(), TeamMember{}) {
		t.Error("TeamMember getType not returning empty team member")
	}
}
