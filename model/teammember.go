package model

/*
TeamMember represents a human individual that is currently employed by the organization.
TeamMembers must have a Name and Title, and a unique ID. TeamMembers may optionally possess
a set of Skills (SkillSet), as well as a set of Skills they wish to obtain (WishList)
 */
type TeamMember struct {
	ID       string
	Name     string
	Title    string
	SkillSet []Skill
	WishList []Skill
}

/*
NewTeamMember is a constructor for the TeamMember type. Returns a new instance of TeamMember,
initialized to the specified ID, Name, and Title, and with an empty SkillSet and WishList.
 */
func NewTeamMember(id, name, title string) TeamMember {
	return TeamMember{
		ID:        id,
		Name:      name,
		Title:     title,
		SkillSet:  []Skill{},
		WishList:  []Skill{},
	}
}

/*
NewTeamMember is a constructor for the TeamMember type. Returns a new instance of TeamMember,
initialized to the specified ID, Name, Title, SkillSet, and WishList.
 */
func NewTeamMemberWithSkills(id, name, title string, skillList, wishList []Skill) TeamMember {
	return TeamMember{
		ID:        id,
		Name:      name,
		Title:     title,
		SkillSet: skillList,
		WishList:  wishList,
	}
}

/*
AddSkill adds a new Skill to the TeamMember's SkillSet.
 */
func (t *TeamMember) AddSkill(skill Skill) {
	t.SkillSet = append(t.SkillSet, skill)
}

/*
AddToWishList adds a new Skill to the TeamMember's WishList.
 */
func (t *TeamMember) AddToWishList(skill Skill) {
	t.WishList = append(t.WishList, skill)
}

// GetType returns an interface{} with an underlying concrete type of TeamMember{}.
func (t *TeamMember) GetType() interface{} {
	return TeamMember{}
}