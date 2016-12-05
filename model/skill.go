package model

type Skill struct {
	Name      string
	SkillType string
}

type SkillType interface {
	getName() string
}

var skillTypeList = []string{"scripted", "compiled", "orchestration", "database"}

func NewSkill(name, skillType string) Skill {
	return Skill{
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
