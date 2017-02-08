package model

type TMSkill struct {
	ID           string `json:"id"`
	SkillID      string `json:"skill_id"`
	TeamMemberID string `json:"team_member_id"`
	WishList     bool   `json:"wish_list"`
	Proficiency  int    `json:"proficiency"`
}

type TMSkillDTO struct {
	TMSkill
	SkillName      string `json:"skill_name"`
	TeamMemberName string `json:"team_member_name"`
}

/*
NewTMSkillDefaults returns a new instance of TMSkill, with defaults for WishList
(false) and Proficiency (0).
*/
func NewTMSkillDefaults(id, skillID, teamMemberID string) TMSkill {
	return TMSkill{
		ID:           id,
		SkillID:      skillID,
		TeamMemberID: teamMemberID,
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
func NewTMSkillSetDefaults(id, skillID, teamMemberID string, wishList bool,
	proficiency int) TMSkill {
	if proficiency > 5 {
		proficiency = 5
	}
	if proficiency < 0 {
		proficiency = 0
	}
	return TMSkill{
		ID:           id,
		SkillID:      skillID,
		TeamMemberID: teamMemberID,
		WishList:     wishList,
		Proficiency:  proficiency,
	}
}

/*
NewTMSkillDTO returns a new TMSkillDTO for the TMSkill it is called on, using
the specified skillName and teamMemberName.
*/
func (t TMSkill) NewTMSkillDTO(skillName, teamMemberName string) TMSkillDTO {
	return TMSkillDTO{
		TMSkill:        t,
		SkillName:      skillName,
		TeamMemberName: teamMemberName,
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
