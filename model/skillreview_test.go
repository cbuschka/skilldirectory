package model

import (
	"reflect"
	"testing"
)

func TestNewSkillReview(t *testing.T) {
	skillReview1 := NewSkillReview(1, 2, 3, "body", true)
	skillReview2 := SkillReview{
		SkillID:      2,
		TeamMemberID: 3,
		Body:         "body",
		Positive:     true,
	}
	skillReview2.ID = 1
	// Verify that all of srOne and srTwo's fields are equal
	if !reflect.DeepEqual(skillReview1, skillReview2) {
		t.Errorf("constructor newSkillReview() produced incorrect SkillReview.")
	}
}

func TestGetSkillReviewType(t *testing.T) {
	s := NewSkillReview(1, 2, 3, "", true)
	if !reflect.DeepEqual(s.GetType(), SkillReview{}) {
		t.Error("SkillReview getType not returning empty skill review")
	}
}

func TestSkillReviewGetID(t *testing.T) {
	s := NewSkillReview(1, 0, 0, "", true)
	if s.GetID() != 1 {
		t.Error("GetID Failed")
	}
}

func TestQuerySkillReview(t *testing.T) {
	one := QuerySkillReview(1)
	two := SkillReview{}
	two.ID = 1
	if !reflect.DeepEqual(one, two) {
		t.Errorf("One: %v doesn't match Two: %v", one, two)
	}
}
