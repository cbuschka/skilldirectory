package gormmodel

import "github.com/jinzhu/gorm"

/*
TeamMember represents a human individual that is currently employed by the
organization. TeamMembers must have a Name and Title, and a unique ID.
TeamMembers may optionally possess a set of Skills (SkillSet), as well as a
set of Skills they wish to obtain (WishList).
*/
type TeamMember struct {
	gorm.Model
	Name  string `json:"name"`
	Title string `json:"title"`
}

// TeamMemberDTO is the transfer object that includes TMSkills that a TeamMember has
type TeamMemberDTO struct {
	TeamMember
	TMSkillDTOs []TMSkillDTO `json:"tm_skills"`
}

/*
NewTeamMember is a constructor for the TeamMember type. Returns a new instance
of TeamMember, initialized to the specified ID, Name, and Title.
*/
func NewTeamMember(id uint, name, title string) TeamMember {
	tm := TeamMember{
		Name:  name,
		Title: title,
	}
	tm.ID = id
	return tm
}

/*
NewTeamMemberDTO returns a new instance of TeamMemberDTO for the TeamMember
it is called on, using the specified []TMSkillDTO.
*/
func (t TeamMember) NewTeamMemberDTO(tmSkillDTOs []TMSkillDTO) TeamMemberDTO {
	return TeamMemberDTO{
		TeamMember:  t,
		TMSkillDTOs: tmSkillDTOs,
	}
}

/*
GetType returns an interface{} with an underlying concrete type of TeamMember{}.
*/
func (t TeamMember) GetType() interface{} {
	return TeamMember{}
}

func (t TeamMember) GetID() uint {
	return t.ID
}

func QueryTeamMember(id uint) TeamMember {
	var tm TeamMember
	tm.ID = id
	return tm
}
