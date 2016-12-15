package model

import "fmt"

type Skill struct {
	Id        string
	Name      string
	SkillType string
	Webpage   Link
	Blogs     []Link
	Tutorials []Link
}

func (s *Skill) AddLink(link Link, linkType string) error {
	if !IsValidLinkType(linkType) {
		return fmt.Errorf("The specified LinkType: \"%s\" is not LinkType.", linkType)
	}

	switch linkType {
	case WebpageLinkType:
		s.Webpage = link
	case BlogLinkType:
		s.Blogs = append(s.Blogs, link)
	case TutorialLinkType:
		s.Tutorials = append(s.Tutorials, link)
	}
	return nil
}

const (
	ScriptedSkillType      = "scripted"
	CompiledSkillType      = "compiled"
	OrchestrationSkillType = "orchestration"
	DatabaseSkillType      = "database"
)

func NewSkill(id, name, skillType string) Skill {
	return Skill{
		Id:        id,
		Name:      name,
		SkillType: skillType,
		Webpage:   Link{},
		Blogs:     []Link{},
		Tutorials: []Link{}}
}

func NewSkillWithLinks(id, name, skillType string,
	webpage Link, blogs, tutorials []Link) Skill {
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
