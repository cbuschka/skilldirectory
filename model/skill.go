package model

type Skill struct {
	Id	string
	Name      string
	SkillType string
}

var skillTypeList = []string{"scripted", "compiled", "orchestration", "database"}

func NewSkill(id, name, skillType string) Skill {
	return Skill{
		Id: id,
		Name:      name,
		SkillType: skillType,
	}
}

func IsValidSkillType(skillType string) bool {
	for _, s := range skillTypeList {
		if skillType == s {
			return true
		}
	}
	return false
}
