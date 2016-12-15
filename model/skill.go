package model

type Skill struct {
	Id        string
	Name      string
	SkillType string
	Webpage   Link
	Blogs     []Link
	Tutorials []Link
}

const (
	ScriptedSkillType      = "scripted"
	CompiledSkillType      = "compiled"
	OrchestrationSkillType = "orchestration"
	DatabaseSkillType      = "database"
)

func NewSkill(id, name, skillType string,
	webpage Link,
	blogs, tutorials []Link) Skill {
	return Skill{
		Id:        id,
		Name:      name,
		SkillType: skillType,
		Webpage:   webpage,
		Blogs:     blogs,
		Tutorials: tutorials,
	}
}

func IsValidSkillType(skillType string) bool {
	switch skillType {
	case
		ScriptedSkillType,
		CompiledSkillType,
		OrchestrationSkillType,
		DatabaseSkillType:
		return true
	}
	return false
}

func (s Skill) GetType() interface{} {
	return Skill{}
}

//
// func(s Skill)GetSlice()[]interface{} {
// 	return []Skill{}
// }
