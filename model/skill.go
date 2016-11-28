package model

type Skill struct {
	Name      string
	SkillType string //language, devops_tool, framework
}

func NewSkill(name, skillType string) Skill {
	return Skill{
		Name:      name,
		SkillType: skillType,
	}
}
