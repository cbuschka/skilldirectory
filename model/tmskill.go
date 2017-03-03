package model

import "github.com/jinzhu/gorm"

// TMSkill has a many-to-one relationship to Skills and TeamMembers
type TMSkill struct {
	gorm.Model
	SkillID      uint `gorm:"index" json:"skill_id"`
	TeamMemberID uint `gorm:"index" json:"team_member_id"`
	Proficiency  uint `json:"proficiency"`
	TeamMember   TeamMember
	Skill        Skill
}

/*
NewTMSkillDefaults returns a new instance of TMSkill, with defaults for WishList
(false) and Proficiency (0).
*/
func NewTMSkillDefaults(id, skillID, teamMemberID uint) TMSkill {
	tmSkill := TMSkill{
		SkillID:      skillID,
		TeamMemberID: teamMemberID,
		Proficiency:  0,
	}
	tmSkill.ID = id
	return tmSkill
}

/*
NewTMSkillDefaultsSetDefaults returns a new instance of TMSkill
*/
func NewTMSkillSetDefaults(id, skillID, teamMemberID, profiency uint) TMSkill {
	tmSkill := TMSkill{
		SkillID:      skillID,
		TeamMemberID: teamMemberID,
		Proficiency:  profiency,
	}
	tmSkill.ID = id
	return tmSkill
}

/*
SetProficiency sets the Proficiency field of the TMSkill instance to the
specified proficiency. The specified proficiency must be in the range of 0-5. If
a value is passed in outside of this range, it is clipped to 0 if it's below 0,
or 5 if it's above 5.
*/
func (t *TMSkill) SetProficiency(proficiency uint) {
	if proficiency > 5 {
		proficiency = 5
	}
	t.Proficiency = uint(proficiency)
}

/*
GetProficiencyString returns a string representation of the TMSkill's Proficiency level.
*/
func (t *TMSkill) GetProficiencyString() string {
	switch t.Proficiency {
	case 0:
		return "Not Applicable"
	case 1:
		return "Fundamentally Aware"
	case 2:
		return "Novice"
	case 3:
		return "Intermediate"
	case 4:
		return "Advanced"
	case 5:
		return "Expert"
	default:
		return "No String Representation Available" // outside range 0-5
	}
}

// GetType returns an interface{} with an underlying concrete type of TMSkill{}.
func (t TMSkill) GetType() interface{} {
	return TMSkill{}
}

func (t TMSkill) GetID() uint {
	return t.ID
}

func QueryTMSKill(id uint) TMSkill {
	var tmskill TMSkill
	tmskill.ID = id
	return tmskill
}
