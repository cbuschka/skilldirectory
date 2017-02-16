package model

import "github.com/jinzhu/gorm"

/*
TeamMember represents a human individual that is currently employed by the
organization. TeamMembers must have a Name and Title, and a unique ID.
TeamMembers may optionally possess a set of Skills (SkillSet), as well as a
set of Skills they wish to obtain (WishList).
*/
type TeamMember struct {
	gorm.Model

	Name  		string 		`json:"name"`
	Title 		string 		`json:"title"`

	TMSkills	[]TMSkill	`json:"teammember_skills"`
}

func NewTeamMember(name, title string) TeamMember {
	return TeamMember{
		Name:  name,
		Title: title,
	}
}

func (t TeamMember) GetType() interface{} {
	return TeamMember{}
}
