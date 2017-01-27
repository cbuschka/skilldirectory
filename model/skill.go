package model

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
type Skill struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SkillType string `json:"skill_type"`
}

// NewSkill returns a new Skill object with specified params
func NewSkill(id, name, skillType string) Skill {
	return Skill{
		ID:        id,
		Name:      name,
		SkillType: skillType,
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

// SkillDTO holds additional content relevant to a Skill
type SkillDTO struct {
	Skill
	Links []Link    `json:"links"`
	Icon  SkillIcon `json:"icon"`
}

// NewSkillDTO returns a new SkillDTO object using receiver and specified params
func (s Skill) NewSkillDTO(links []Link, icon SkillIcon) SkillDTO {
	return SkillDTO{
		Skill: s,
		Links: links,
		Icon:  icon,
	}
}

// AddLink adds a new link to the receiver's Links slice
func (s *SkillDTO) AddLink(link Link) {
	s.Links = append(s.Links, link)
}
