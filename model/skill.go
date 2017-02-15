package model

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
	 
*/
type Skill struct {
	gorm.Model

	Name     			string 				`json:"name"`
	SkillType			string 				`json:"skill_type"`

	IconURL				string 				`json:"icon_url"`

	Links					[]Link 				`json:"links"`
	SkillReviews	[]SkillReview	`json:"skills"`

	TeamMembers		[]TeamMember	`gorm:"many2many:teammember_skills" json:"teamembers"`
}

// NewSkill returns a new Skill object with specified params
func NewSkill(name, skillType, iconURL string) Skill {
	return Skill{
		Name:      name,
		SkillType: skillType,
		IconURL:   iconURL,
	}
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
