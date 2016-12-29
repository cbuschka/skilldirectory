package model

import "testing"

func TestNewSkillReview(t *testing.T) {
	srOne := NewSkillReview("1234", "2345", "3456", "blahblahblah", "1234", true)
	srTwo := SkillReview{
		ID:           "1234",
		SkillID:      "2345",
		TeamMemberID: "3456",
		Body:         "blahblahblah",
		Timestamp:    "1234",
		Positive:     true,
	}
	// Verify that all of srOne and srTwo's fields are equal
	if srOne != srTwo {
		t.Errorf("constructor newSkillReview() produced incorrect SkillReview.")
	}
}