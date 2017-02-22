package gormmodel

// func TestNewSkillReview(t *testing.T) {
// 	srDTOOne := NewSkillReview("1234", "2345", "3456", "blahblahblah", "1234",
// 		true).NewSkillReviewDTO("Coding", "Javatar")
// 	srDTOTwo := SkillReviewDTO{
// 		SkillReview: SkillReview{
// 			ID:           "1234",
// 			SkillID:      "2345",
// 			TeamMemberID: "3456",
// 			Body:         "blahblahblah",
// 			Timestamp:    "1234",
// 			Positive:     true,
// 		},
// 		SkillName:      "Coding",
// 		TeamMemberName: "Javatar",
// 	}
// 	// Verify that all of srOne and srTwo's fields are equal
// 	if !reflect.DeepEqual(srDTOOne, srDTOTwo) {
// 		t.Errorf("constructor newSkillReview() produced incorrect SkillReview.")
// 	}
// }
//
// func TestGetSkillReviewType(t *testing.T) {
// 	s := NewSkillReview("", "", "", "", "", true)
// 	if !reflect.DeepEqual(s.GetType(), SkillReview{}) {
// 		t.Error("SkillReview getType not returning empty skill review")
// 	}
// }
