package gormmodel

import "github.com/jinzhu/gorm"

/*
SkillReview represents a review of a particular Skill. SkillReviews can be
positive or negative (determined by Positive flag). Each review must be linked
to a specific Skill and TeamMember, and must also contain a date and body
(substance of the review).
*/
type SkillReview struct {
	gorm.Model
	Body         string `json:"body"`
	Positive     bool   `json:"positive"`
	SkillID      uint   `gorm:"index" json:"skill_id"`
	TeamMemberID uint   `gorm:"index" json:"team_member_id"`
	TeamMember   TeamMember
	Skills       Skill
}

/*
SkillReviewDTO holds a SkillReview, as well as the names of the Skill and
TeamMember that the SkillReview holds the IDs of.
*/

/*
NewSkillReview returns a new instance of SkillReview. All fields must be specified.
*/
func NewSkillReview(id, skillID, teamMemberID uint, body string, positive bool) SkillReview {
	skillReview := SkillReview{
		SkillID:      skillID,
		TeamMemberID: teamMemberID,
		Body:         body,
		Positive:     positive,
	}
	skillReview.ID = id
	return skillReview
}

// GetType returns an interface{} with an underlying concrete type of SkillReview
func (s SkillReview) GetType() interface{} {
	return SkillReview{}
}

func (s SkillReview) GetID() uint {
	return s.ID
}
