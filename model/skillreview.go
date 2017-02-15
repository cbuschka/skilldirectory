package model

import "github.com/jinzhu/gorm"
/*
SkillReview represents a review of a particular Skill. SkillReviews can be
positive or negative (determined by Positive flag). Each review must be linked
to a specific Skill and TeamMember, and must also contain a date and body
(substance of the review).
*/
type SkillReview struct {
	gorm.Model
	Body        	string	`json:"body"`
	Positive    	bool		`json:"positive"`

	SkillID				uint		`gorm:"index"`
	TeamMemberID	uint		`gorm:"index"`
}

/*
NewSkillReview returns a new instance of SkillReview. All fields must be specified.
*/
func NewSkillReview(body string, positive bool) SkillReview {
	return SkillReview{
		Body:         body,
		Positive:     positive,
	}
}

// GetType returns an interface{} with an underlying concrete type of SkillReview
func (s SkillReview) GetType() interface{} {
	return SkillReview{}
}
