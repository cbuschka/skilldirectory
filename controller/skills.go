package controller

import (
	"encoding/json"
	"net/url"
	"path"
	"skilldirectory/model"
)

type SkillsController struct {
	BaseController
}

func (c SkillsController) Get() {
	c.performGet()
}
func (c SkillsController) Post() {}

func (c SkillsController) Delete() {}

func (c SkillsController) Put() {}

func (c SkillsController) performGet() error {
	path := checkForId(c.r.URL)
	if path == "" {
		return c.getAllSkills()
	}
	return c.getSkill(path)
}

func (c *SkillsController) getAllSkills() error {
	skills, err := c.session.ReadAll("skills/", model.Skill{})
	if err != nil {
		return err
	}
	b, err := json.Marshal(skills)
	c.w.Write(b)
	return err
}

func (c *SkillsController) getSkill(id string) error {
	skill, err := c.loadSkill(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(skill)
	c.w.Write(b)
	return err
}

func (c *SkillsController) loadSkill(id string) (*model.Skill, error) {
	skill := model.Skill{}
	err := c.session.Read(id, &skill)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func checkForId(url *url.URL) string {
	base := path.Base(url.RequestURI())
	if base == "skills" {
		return ""
	}
	return base
}
