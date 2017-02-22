package gormmodel

/*
SkillReview represents a review of a particular Skill. SkillReviews can be
positive or negative (determined by Positive flag). Each review must be linked
to a specific Skill and TeamMember, and must also contain a date and body
(substance of the review).
*/
type SkillReview struct {
	ID           string `json:"id"`
	SkillID      string `json:"skill_id"`
	TeamMemberID string `json:"team_member_id"`
	Body         string `json:"body"`
	Timestamp    string `json:"timestamp"`
	Positive     bool   `json:"positive"`
}

/*
SkillReviewDTO holds a SkillReview, as well as the names of the Skill and
TeamMember that the SkillReview holds the IDs of.
*/
type SkillReviewDTO struct {
	SkillReview
	SkillName      string `json:"skill_name"`
	TeamMemberName string `json:"team_member_name"`
}

/*
NewSkillReview returns a new instance of SkillReview. All fields must be specified.
*/
func NewSkillReview(id, skillID, teamMemberID, body, timestamp string, positive bool) SkillReview {
	return SkillReview{
		ID:           id,
		SkillID:      skillID,
		TeamMemberID: teamMemberID,
		Body:         body,
		Timestamp:    timestamp,
		Positive:     positive,
	}
}

/*
NewSkillReviewDTO returns a new SkillReviewDTO for the SkillReview it is called
on, using the specified skillName and teamMemberName.
*/
func (s SkillReview) NewSkillReviewDTO(skillName,
	teamMemberName string) SkillReviewDTO {
	return SkillReviewDTO{
		SkillReview:    s,
		SkillName:      skillName,
		TeamMemberName: teamMemberName,
	}
}

// GetType returns an interface{} with an underlying concrete type of SkillReview
func (s SkillReview) GetType() interface{} {
	return SkillReview{}
}
