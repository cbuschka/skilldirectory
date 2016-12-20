package model

type TMSkill struct {
	ID           string `json:"ID"`
	SkillID      string `json:"skill_id`
	TeamMemberID string `json:"team_member_id`
	WishList     bool   `json:"wish_list"`
	SkillSet     bool   `json:"skill_set"`
	Proficiency  int    `json:"proficiency"`
}

/*
NewTMSkillDefaults returns a new instance of TMSkill, with defaults for WishList (false), 
SkillSet (false), and Proficiency (0).
 */
func NewTMSkillDefaults(id, skillID, teamMemberID string) TMSkill {
	return TMSkill{
		ID:           id,
		SkillID:      skillID,
		TeamMemberID: teamMemberID,
		WishList:     false,
		SkillSet:     false,
		Proficiency:  0,
	}
}

/*
NewTMSkillSetDefaults returns a new instance of TMSkill, with all fields specified by the caller.
The proficiency field must be in the range of 0-5. If a value is passed in outside of this range, it
is clipped to 0 if it's below 0, or 5 if it's above 5.
 */
func NewTMSkillSetDefaults(id, skillID, teamMemberID string, wishList, skillSet bool, proficiency int) TMSkill {
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
		SkillSet:     skillSet,
		Proficiency:  proficiency,
	}
}

/*
setProficiency sets the Proficiency field of the TMSkill instance to the specified proficiency.
The specified proficiency must be in the range of 0-5. If a value is passed in outside of this range, it
is clipped to 0 if it's below 0, or 5 if it's above 5.
 */
func (tmSkill *TMSkill) SetProficiency(proficiency int) {
	if proficiency > 5 {
		proficiency = 5
	}
	if proficiency < 0 {
		proficiency = 0
	}
	tmSkill.Proficiency = proficiency
}

// GetType returns an interface{} with an underlying concrete type of TMSkill{}.
func (s TMSkill) GetType() interface{} {
	return TMSkill{}
}