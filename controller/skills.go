package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"skilldirectory/model"

	"skilldirectory/errors"

	"fmt"

	uuid "github.com/satori/go.uuid"
)

type SkillsController struct {
	*BaseController
}

func (c SkillsController) Base() *BaseController {
	return c.BaseController
}
func (c SkillsController) Get() error {
	return c.performGet()
}
func (c SkillsController) Post() error {
	return c.addSkill()
}

func (c SkillsController) Delete() error {
	return c.removeSkill()
}

func (c SkillsController) Put() error {
	return nil
}

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
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No Skill Exists with Specified ID: " + id,
		}
	}
	return &skill, nil
}

func (c *SkillsController) removeSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID := checkForId(c.r.URL)
	if skillID == "" {
		return &errors.MissingSkillIDError{
			ErrorMsg: "No Skill ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete(skillID)
	if err != nil {
		return &errors.NoSuchIDError{
			ErrorMsg: "No Skill Exists with Specified ID: " + skillID,
		}
	}

	log.Printf("Skill Deleted with ID: %s", skillID)
	return nil
}

func (c *SkillsController) addSkill() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	skill := model.Skill{}
	err := json.Unmarshal(body, &skill)
	if err != nil {
		return &errors.MarshalingError{
			ErrorMsg: "Invalid JSON body in request:\n\t" + fmt.Sprint(body),
		}
	}
	if !model.IsValidSkillType(skill.SkillType) {
		return &errors.InvalidSkillTypeError{
			ErrorMsg: "Invalid Skill Type: %s" + skill.SkillType,
		}
	}
	skill.Id = uuid.NewV1().String()
	err = c.session.Save(skill.Id, skill)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
	log.Printf("Saved skill: %s", skill.Name)
	return nil
}
