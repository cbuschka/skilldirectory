package model

import "github.com/jinzhu/gorm"

// TMSkill has a many-to-one relationship to Skills and TeamMembers
type TMSkill struct {
	gorm.Model
	WishList     	bool   			`json:"wish_list"`
	Proficiency  	int    			`json:"proficiency"`

	SkillID				uint				`gorm:"index"`
	Skill					Skill				`json:"skill"`
	TeamMemberID	uint				`gorm:"index"`
	TeamMember 		TeamMember	`json:"team_member"`
}

/*
NewTMSkillDefaults returns a new instance of TMSkill, with defaults for WishList
(false) and Proficiency (0).
*/
func NewTMSkillDefaults() TMSkill {
	return TMSkill{
		WishList:     false,
		Proficiency:  0,
	}
}

/*
NewTMSkillSetDefaults returns a new instance of TMSkill, with all fields
specified by the caller. The proficiency field must be in the range of 0-5. If a
value is passed in outside of this range, it is clipped to 0 if it's below 0, or
5 if it's above 5.
*/
func NewTMSkillSetDefaults(wishList bool, proficiency int) TMSkill {
	if proficiency > 5 {
		proficiency = 5
	}
	if proficiency < 0 {
		proficiency = 0
	}
	return TMSkill{
		WishList:     wishList,
		Proficiency:  proficiency,
	}
}

/*
SetProficiency sets the Proficiency field of the TMSkill instance to the
specified proficiency. The specified proficiency must be in the range of 0-5. If
a value is passed in outside of this range, it is clipped to 0 if it's below 0,
or 5 if it's above 5.
*/
func (t *TMSkill) SetProficiency(proficiency int) {
	if proficiency > 5 {
		proficiency = 5
	}
	if proficiency < 0 {
		proficiency = 0
	}
	t.Proficiency = proficiency
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
