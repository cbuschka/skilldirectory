package model

import (
	"reflect"
	"testing"
)

func TestNewSkillReview(t *testing.T) {
	srOne := NewSkillReview("blahblahblah", true)
	srTwo := SkillReview {
			Body:         "blahblahblah",
			Positive:     true,
	}
	// Verify that all of srOne and srTwo's fields are equal
	if !reflect.DeepEqual(srOne, srTwo) {
		t.Errorf("constructor newSkillReview() produced incorrect SkillReview.")
	}
}

func TestGetSkillReviewType(t *testing.T) {
	s := NewSkillReview("", true)
	if !reflect.DeepEqual(s.GetType(), SkillReview{}) {
		t.Error("SkillReview getType not returning empty skill review")
	}
}
