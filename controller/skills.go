package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"skilldirectory/model"

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
		return nil, err
	}
	return &skill, nil
}

// Removes the skill with the ID at the end of the specified request's URL.
// Returns non-nil error if the request's URL contains no ID, or if no skills
// exist with that ID.
func (c *SkillsController) removeSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID := checkForId(c.r.URL)
	if skillID == "" {
		return fmt.Errorf("No Skill ID Specified in Request URL: %s", c.r.URL)
	}

	// Remove the skill with the specified ID from the skills
	// database/repository; return error if no skills have that ID
	err := c.session.Delete(skillID)
	if err != nil {
		return fmt.Errorf("No Skill Exists with Specified ID: %s", skillID)
	}

	// The skill was successfully deleted!
	log.Printf("Skill Deleted with ID: %s", skillID)
	return nil
}

func (c *SkillsController) addSkill() error {
	fmt.Println("Adding skill")
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	skill := model.Skill{}
	err := json.Unmarshal(body, &skill)
	if err != nil {
		return err
	}
	if !model.IsValidSkillType(skill.SkillType) {
		return fmt.Errorf("Invalid Skill Type: %s", skill.SkillType)
	}
	skill.Id = uuid.NewV1().String()
	err = c.session.Save(skill.Id, skill)
	if err != nil {
		return err
	}
	log.Printf("New skill saved")
	return nil
}
