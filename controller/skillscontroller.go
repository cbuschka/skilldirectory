package controller

import (
	"encoding/json"
	"io/ioutil"
	"skilldirectory/data"
	"skilldirectory/model"

	"skilldirectory/errors"
	"skilldirectory/gormmodel"
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

func (c SkillsController) Options() error {
	c.w.Header().Set("Access-Control-Allow-Headers", GetDefaultHeaders())
	c.w.Header().Set("Access-Control-Allow-Methods", GetDefaultMethods())
	return nil
}

func (c SkillsController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllSkills()
	}

	skillID, err := util.PathToID(c.r.URL)
	if err != nil {
		return fmt.Errorf("ID must be an uint")
	}
	return c.getSkill(skillID)
}

func (c *SkillsController) getAllSkills() error {
	var err error
	// filter := c.r.URL.Query().Get("skilltype")
	// var opts data.CassandraQueryOptions

	// // Add approved query filters here
	// if filter != "" {
	// 	opts = data.NewCassandraQueryOptions("skilltype", filter, false)
	// }
	// skills, err = c.session.FilteredReadAll("skills", opts, model.SkillDTO{})

	var skills []model.Skill
	err = c.find(&skills)

	if err != nil {
		return err
	}

	b, err := json.Marshal(skills)
	c.w.Write(b)
	return err
}

func (c *SkillsController) getSkill(id uint) error {
	skill := gormmodel.QuerySkill(id)
	err := c.first(&skill)
	if err != nil {
		return err
	}
	b, err := json.Marshal(skill)
	c.w.Write(b)
	return err
}

func (c *SkillsController) removeSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID := util.CheckForID(c.r.URL)
	if skillID == "" {
		return errors.MissingIDError(fmt.Errorf("no Skill ID in request URL"))
	}
	skill := gormmodel.Skill{}
	id, err := util.PathToID(c.r.URL)
	if err != nil {
		return err
	}
	skill.ID = id
	err = c.delete(skill)

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

	c.Printf("Adding to Postgres")
	skill := gormmodel.Skill{}
	err := json.Unmarshal(body, &skill)
	if err != nil {
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
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

	if err != nil {
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
	}

	err = c.create(&skill)
	if err != nil {
		return err
	}

	// Return object JSON as response
	b, err := json.Marshal(skill)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	c.Printf("Saved skill: %s", skill.Name)
	return nil
}

/*
validatePOSTBody() accepts a model.Skill pointer. It can be used to verify the
validity of the state of a Skill initialized via unmarshaled JSON. Ensures that the
passed-in Skill contains a key-value pair for "Name" and for "SkillType"
fields. Returns nil error if it does, IncompletePOSTBodyError error if not.
*/
func (c *SkillsController) validatePOSTBody(skill *gormmodel.Skill) error {
	if skill.Name == "" || skill.SkillType == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"A Skill must be a JSON object and must contain values for "+
				"%q and %q fields.", "name", "skill_type"))
	}
	return nil
}
