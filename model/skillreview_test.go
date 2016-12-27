package model

import "testing"

func TestNewSkillReview(t *testing.T) {
	srOne := NewSkillReview("1234", "2345", "3456", "blahblahblah", "12/27/2016", true)
	srTwo := SkillReview{
		ID:           "1234",
		SkillID:      "2345",
		TeamMemberID: "3456",
		Body:         "blahblahblah",
		Date:         "12/27/2016",
		Positive:     true,
	}
	// Verify that all of srOne and srTwo's fields are equal
	if srOne != srTwo {
		t.Errorf("constructor newSkillReview() produced incorrect SkillReview.")
	}
}