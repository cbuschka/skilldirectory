package model

/*
Skill represents a particular skill that can be had by a human individual.
Each Skill has a Name, SkillType, and a unique ID:

 * The Name should appropriately identify the skill, such as "Java", "SQL", "Go", or "Baking Cookies".

 * The SkillType must be one of the predetermined SkillTypes contained within model/skills.go as
   constants (e.g. "models.ScriptedSkillType" or "DatabaseSkillType").

 * The ID can be any desired string value, but ought to be unique, so that it can
   be used to identify the skill should it be stored in a database with other Skills.
*/
type Skill struct {
	ID        string
	Name      string
	SkillType string
}

const (
	// e.g. writing Python or Bash scripts
	ScriptedSkillType = "scripted"
	// e.g. writing Java or C++ code
	CompiledSkillType = "compiled"
	// Ability to integrate multiple services to automate a process and provide a single, unified service.
	OrchestrationSkillType = "orchestration"
	// e.g. SQL or JDBC knowledge and aptitude
	DatabaseSkillType = "database"
)

// NewSkill() creates and returns a new instance of model.Skill
func NewSkill(id, name, skillType string) Skill {
	return Skill{
		ID:        id,
		Name:      name,
		SkillType: skillType,
	}
}

// IsValidSkillType() returns true if the passed-in string is a valid SkillType, false if not.
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

// GetType() returns an interface{} with an underlying concrete type of Skill{}.
func (s Skill) GetType() interface{} {
	return Skill{}
}
