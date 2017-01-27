package controller

import (
	"encoding/json"
	"io/ioutil"
	"skilldirectory/data"
	"skilldirectory/model"

	"skilldirectory/errors"
	"skilldirectory/util"

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
	err := c.session.Read("skills", id, data.CassandraQueryOptions{}, &skill)
	if err != nil {
		return nil, errors.NoSuchIDError(fmt.Errorf(
			"no Skill exists with specified ID: %s", id))
	}
	skillDTO, _ := c.addLinks(skill)
	c.addIcon(&skillDTO)
	return &skillDTO, nil
}

func (c *SkillsController) addLinks(skill model.Skill) (model.SkillDTO, error) {
	skillDTO := model.SkillDTO{}
	linksInterface, err := c.session.FilteredReadAll("links",
		data.NewCassandraQueryOptions("skill_id", skill.ID, true), model.Link{})
	if err != nil {
		c.Print(err)
		return skillDTO, err
	}
	linksRaw, err := json.Marshal(linksInterface)
	if err != nil {
		return skillDTO, nil
	}
	links := &[]model.Link{}
	err = json.Unmarshal(linksRaw, links)
	if err != nil {
		c.Print(err)
	}
	skillDTO = skill.NewSkillDTO(*links, model.SkillIcon{})
	return skillDTO, nil
}

func (c *SkillsController) addIcon(skillDTO *model.SkillDTO) {
	skillIcon := model.SkillIcon{}
	c.session.Read("skillicons", "",
		data.NewCassandraQueryOptions("skill_id", skillDTO.Skill.ID, true), &skillIcon)
	skillDTO.Icon = skillIcon
}

func (c *SkillsController) removeSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID := util.CheckForID(c.r.URL)
	if skillID == "" {
		return errors.MissingIDError(fmt.Errorf("no Skill ID in request URL"))
	}
	err1 := c.removeSkillChildren(skillID)
	if err1 != nil {
		c.Printf("removingSkillChildren: %v", err1)

	}

	err := c.session.Delete("skills", skillID, data.CassandraQueryOptions{})

	if err != nil {
		c.Printf("removeSkill() failed for the following reason:\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no Skill exists with specified ID: %s", skillID))
	}

	c.Printf("Skill Deleted with ID: %s", skillID)
	return nil
}

func (c *SkillsController) removeSkillChildren(skillID string) error {
	return c.session.Delete("links", "", data.NewCassandraQueryOptions("skill_ID", skillID, true))

}

func (c *SkillsController) addSkill() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	skill := model.Skill{}
	err := json.Unmarshal(body, &skill)
	if err != nil {
		return errors.MarshalingError(err)
	}

	err = c.validatePOSTBody(&skill)
	if err != nil {
		c.Debugf("Invalid Post: Name: %s, ID: %s", skill.Name, skill.SkillType)
		return err // Will be of errors.IncompletePOSTBodyError type
	}

	if !model.IsValidSkillType(skill.SkillType) {
		return errors.InvalidSkillTypeError(fmt.Errorf(
			"invalid Skill type: %s", skill.SkillType))
	}

	skill.ID = util.NewID()
	err = c.session.Save("skills", skill.ID, skill)
	if err != nil {
		return errors.SavingError(err)
	}
	c.Printf("Saved skill: %s", skill.Name)
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
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"The JSON in a POST Request for new Skill must contain values for "+
				"%q and %q fields.", "name", "skill_type"))
	}
	return nil
}
