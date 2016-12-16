package model

/*
TeamMember represents a human individual that is currently employed by the organization.
TeamMembers must have a Name and Title, and a unique ID. TeamMembers may optionally possess
a set of Skills (SkillSet), as well as a set of Skills they wish to obtain (WishList)
*/
type TeamMember struct {
	ID    string
	Name  string
	Title string
}

/*
NewTeamMember is a constructor for the TeamMember type. Returns a new instance of TeamMember,
initialized to the specified ID, Name, and Title.
*/
func NewTeamMember(id, name, title string) TeamMember {
	return TeamMember{
		ID:    id,
		Name:  name,
		Title: title,
	}
}

// GetType returns an interface{} with an underlying concrete type of TeamMember{}.
func (t TeamMember) GetType() interface{} {
	return TeamMember{}
}
