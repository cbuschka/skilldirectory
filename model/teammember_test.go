package model

import "testing"

func TestNewTeamMember(t *testing.T) {
	tmOne := NewTeamMember("1234", "Yogi Bear", "Smarter Than Average Bear")
	tmTwo := TeamMember{
		ID: "1234",
		Name: "Yogi Bear",
		Title: "Smarter Than Average Bear",
	}
	if tmOne != tmTwo {
		t.Error("model.teammember.\"NewTeamMember()\" produced incorrect TeamMember")
	}
}