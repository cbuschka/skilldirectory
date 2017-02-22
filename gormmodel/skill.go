package gormmodel

import "github.com/jinzhu/gorm"

const (
	// ScriptedSkillType indicates a skill like writing Python or Bash scripts
	ScriptedSkillType = "scripted"
	// CompiledSkillType indicates a skill like writing Java or C++ code
	CompiledSkillType = "compiled"
	// OrchestrationSkillType indicates a skill such as the ability to integrate
	// multiple services to automate a process and provide a single, unified
	// service.
	OrchestrationSkillType = "orchestration"
	// DatabaseSkillType indicates knowledge in an area such as SQL or JDBC
	DatabaseSkillType = "database"
)

/*
Skill models a particular skill that can be had by a human individual.
Each Skill has a Name, SkillType, and a unique ID:
 * The Name should appropriately identify the skill, such as "Java", "SQL",
   "Go", or "Baking Cookies".

 * The SkillType must be one of the predetermined SkillTypes contained within
   model/skills.go as

 * The ID can be any desired string value, but ought to be unique, so that it
   can be used to identify the skill should it be stored in a database with
	 other Skills.
*/

type GormInterface interface {
	GetID() uint
}
type Skill struct {
	gorm.Model

	Name      string `json:"name"`
	SkillType string `json:"skill_type"`

	IconURL string `json:"icon_url"`

	Links        []Link
	SkillReviews []SkillReview
	TMSkills     []TMSkill
}

func (s Skill) GetID() uint {
	return s.ID
}

func QuerySkill(id uint) Skill {
	var skill Skill
	skill.ID = id
	return skill
}

// NewSkill returns a new Skill object with specified params
func NewSkill(id uint, name, skillType string) Skill {
	skill := Skill{
		Name:      name,
		SkillType: skillType,
	}
	skill.ID = id
	return skill
}

// GetType satisfies data.ReadAllInterface
func (s Skill) GetType() interface{} {
	return Skill{}
}

// IsValidSkillType returns true if skillType is a valid SkillType, false if not.
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
