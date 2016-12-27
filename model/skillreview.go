package model

/*
SkillReview represents a review of a particular Skill. SkillReviews can be positive
or negative (determined by Positive flag). Each review must be linked to a specific
Skill and TeamMember, and must also contain a date and body (substance of the review).
 */
type SkillReview struct {
	ID           string `json:"id"`;
	SkillID      string `json:"skill_id"`;
	TeamMemberID string `json:"team_member_id"`;
	Body         string `json:"body"`;
	Date         string `json:"date"`;
	Positive     bool   `json:"positive"`;
}

/*
NewSkillReview returns a new instance of SkillReview. All fields must be specified.
 */
func NewSkillReview(id, skillID, teamMemberID, body, date string, positive bool) SkillReview {
	return SkillReview{
		ID: id,
		SkillID: skillID,
		TeamMemberID: teamMemberID,
		Body: body,
		Date: date,
		Positive: positive,
	}
}

// GetType returns an interface{} with an underlying concrete type of SkillReview{}.
func (s SkillReview) GetType() interface{} {
	return SkillReview{}
}
