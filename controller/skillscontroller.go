package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"skilldirectory/data"
	"skilldirectory/model"

	"skilldirectory/errors"
	util "skilldirectory/util"

	"fmt"
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
	return fmt.Errorf("PUT requests not currently supported.")
}

func (c SkillsController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllSkills()
	}
	return c.getSkill(path)
}

func (c *SkillsController) getAllSkills() error {
	var skills []interface{}
	var err error
	filter := c.r.URL.Query().Get("skilltype")
	var opts data.CassandraQueryOptions

	// Add approved query filters here
	if filter != "" {
		opts = data.NewCassandraQueryOptions("skilltype", filter, false)
	}
	skills, err = c.session.FilteredReadAll("skills", opts, model.SkillDTO{})

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

func (c *SkillsController) loadSkill(id string) (*model.SkillDTO, error) {
	skill := model.Skill{}
	err := c.session.Read("skills", id, &skill)
	if err != nil {
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No Skill Exists with Specified ID: " + id,
		}
	}
	skillDTO, _ := c.addLinks(skill)
	return &skillDTO, nil
}

func (c *SkillsController) addLinks(skill model.Skill) (model.SkillDTO, error) {
	skillDTO := model.SkillDTO{}
	linksInterface, err := c.session.FilteredReadAll("links", data.NewCassandraQueryOptions("skill_id", skill.ID, true), model.Link{})
	if err != nil {
		log.Print(err)
		return skillDTO, err
	}
	linksRaw, err := json.Marshal(linksInterface)
	if err != nil {
		return skillDTO, nil
	}
	links := &[]model.Link{}
	err = json.Unmarshal(linksRaw, links)
	if err != nil {
		log.Print(err)
	}
	skillDTO = skill.NewSkillDTO(*links)
	return skillDTO, nil
}

func (c *SkillsController) removeSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID := util.CheckForID(c.r.URL)
	if skillID == "" {
		return &errors.MissingIDError{
			ErrorMsg: "No Skill ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete("skills", skillID)
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
			ErrorMsg: err.Error(),
		}
	}

	err = c.validatePOSTBody(&skill)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError type
	}

	if !model.IsValidSkillType(skill.SkillType) {
		return &errors.InvalidSkillTypeError{
			ErrorMsg: fmt.Sprintf("Invalid Skill Type: %s", skill.SkillType),
		}
	}

	skill.ID = util.NewID()
	err = c.session.Save("skills", skill.ID, skill)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
	log.Printf("Saved skill: %s", skill.Name)
	return nil
}

/*
validatePOSTBody() accepts a model.Skill pointer. It can be used to verify the
validity of the state of a Skill initialized via unmarshaled JSON. Ensures that the
passed-in Skill contains a key-value pair for "Name" and for "SkillType"
fields. Returns nil error if it does, IncompletePOSTBodyError error if not.
*/
func (c *SkillsController) validatePOSTBody(skill *model.Skill) error {
	if skill.Name == "" || skill.SkillType == "" {
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: "POST Request for new Skill must contain values for " +
				"\"Name\" and \"SkillType\" fields.",
		}
	}
	return nil
}
