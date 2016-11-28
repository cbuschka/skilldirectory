package model

import (
	"encoding/json"
	"io/ioutil"
)

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

func (s Skill) Save() error {
	path := pathHelper(s.Name)
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0600)
}

func ReadSkill(name string) (Skill, error) {
	path := pathHelper(name)
	skill := Skill{}
	data, err := ioutil.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &skill)
	}
	return skill, err

}

func pathHelper(name string) string {
	return "skills/" + name + ".txt"
}
