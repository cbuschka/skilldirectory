package controller

import (
	"encoding/json"
	"io/ioutil"

	"skilldirectory/errors"
	"skilldirectory/model"
	"skilldirectory/util"

	"fmt"
)

// SkillsController handles requests for the Skill type
type SkillsController struct {
	*BaseController
}

// Base implemented
func (c SkillsController) Base() *BaseController {
	return c.BaseController
}

// Get implemented
func (c SkillsController) Get() error {
	return c.performGet()
}

// Post implemented
func (c SkillsController) Post() error {
	return c.addSkill()
}

// Delete implemented
func (c SkillsController) Delete() error {
	return c.removeSkill()
}

// Put implemented
func (c SkillsController) Put() error {
	return fmt.Errorf("PUT requests not currently supported")
}

// Options implemented
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
	var skills []model.Skill

	filter := c.r.URL.Query().Get("skilltype")
	// Add approved query filters here
	if filter != "" {
		filterMap := util.NewFilterMap("skill_type", filter)
		err = c.findWhere(&skills, filterMap)
	} else {
		err = c.find(&skills)
	}

	if err != nil {
		return err
	}

	b, err := json.Marshal(skills)
	c.w.Write(b)
	return err
}

func (c *SkillsController) getSkill(id uint) error {
	skill := model.QuerySkill(id)
	err := c.preloadAndFind(&skill, "Links", "SkillReviews")
	if err != nil {
		return err
	}

	c.populateSkillReviews(&skill)
	b, err := json.Marshal(skill)
	c.w.Write(b)
	return err
}

func (c *SkillsController) populateSkillReviews(skill *model.Skill) {
	for i := range skill.SkillReviews {
		review := &skill.SkillReviews[i]
		err := c.preloadAndFind(&review, "TeamMember")
		if err != nil {
			c.Printf("Preload TeamMembers Error: %v", err)
		}
	}
}

func (c *SkillsController) removeSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID, err := util.PathToID(c.r.URL)
	if err != nil {
		return err
	}
	skill := model.QuerySkill(skillID)
	err = c.delete(skill)
	if err != nil {
		c.Printf("removeSkill() failed for the following reason:\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no Skill exists with specified ID: %d", skillID))
	}

	c.Printf("Skill Deleted with ID: %d", skillID)
	return nil
}

func (c *SkillsController) addSkill() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	var skill model.Skill
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

	err = c.create(&skill)
	if err != nil {
		return errors.SavingError(err)
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
func (c *SkillsController) validatePOSTBody(skill *model.Skill) error {
	if skill.Name == "" || skill.SkillType == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"A Skill must be a JSON object and must contain values for "+
				"%q and %q fields.", "name", "skill_type"))
	}
	return nil
}
