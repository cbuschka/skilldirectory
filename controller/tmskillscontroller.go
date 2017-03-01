package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/errors"
	"skilldirectory/gormmodel"
	util "skilldirectory/util"
)

// TMSkillsController handles TMSkills Requests
type TMSkillsController struct {
	*BaseController
}

// Base implemented
func (c TMSkillsController) Base() *BaseController {
	return c.BaseController
}

// Get implemented
func (c TMSkillsController) Get() error {
	return c.performGet()
}

// Post implemented
func (c TMSkillsController) Post() error {
	return c.addTMSkill()
}

// Delete implemented
func (c TMSkillsController) Delete() error {
	return c.removeTMSkill()
}

// Put implemented
func (c TMSkillsController) Put() error {
	return c.updateTMSkill()
}

func (c TMSkillsController) Options() error {
	c.w.Header().Set("Access-Control-Allow-Headers", GetDefaultHeaders())
	c.w.Header().Set("Access-Control-Allow-Methods", "PUT, "+GetDefaultMethods())
	return nil
}

func (c *TMSkillsController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllTMSkills()
	}

	tmSkillID, err := util.StringToID(path)
	if err != nil {
		return err
	}
	return c.getTMSkill(tmSkillID)
}

func (c *TMSkillsController) getAllTMSkills() error {
	var tmSkills []gormmodel.TMSkill
	err := c.find(&tmSkills)
	if err != nil {
		return err
	}
	b, err := json.Marshal(tmSkills)
	c.w.Write(b)
	return err
}

func (c *TMSkillsController) getTMSkill(id uint) error {
	tmSkill := gormmodel.QueryTMSKill(id)
	err := c.first(&tmSkill)
	if err != nil {
		return err
	}
	teamMember := gormmodel.QueryTeamMember(tmSkill.TeamMemberID)
	err = c.first(&teamMember)
	if err != nil {
		return errors.NoSuchIDError(fmt.Errorf(
			"no TeamMember exists with specified ID: %q", tmSkill.TeamMemberID))
	}

	skill := gormmodel.QuerySkill(tmSkill.SkillID)
	err = c.first(&skill)
	if err != nil {
		c.Warnf("Possible invalid id: %v", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no Skill exists with specified ID: %q", tmSkill.SkillID))
	}
	tmSkill.TeamMember = teamMember
	tmSkill.Skill = skill

	b, err := json.Marshal(tmSkill)
	c.w.Write(b)
	return err
}

func (c *TMSkillsController) removeTMSkill() error {
	// Get the ID at end of the request; return error if request contains no ID
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return errors.MissingIDError(fmt.Errorf("no TMSkill ID in request URL"))
	}

	tmSkillID, err := util.StringToID(path)
	if err != nil {
		return err
	}

	tmSkill := gormmodel.QueryTMSKill(tmSkillID)
	err = c.delete(&tmSkill)
	if err != nil {
		c.Printf("removeTMSkill() failed for the following reason:\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no TMSkill exists with specified ID: %q", tmSkillID))
	}

	c.Printf("TMSkill Deleted with ID: %d", tmSkillID)
	return nil
}

// Updates specific TMSkill for PUT requests to "/tmskills/[ID]"
func (c *TMSkillsController) updateTMSkill() error {
	// Get the ID at end of the request; return error if request contains no ID
	tmSkillID := util.CheckForID(c.r.URL)
	if tmSkillID == "" {
		return errors.MissingIDError(fmt.Errorf(
			"must specify a TMSkill ID in PUT request URL"))
	}

	// Store request's body in raw byte slice
	body, err := ioutil.ReadAll(c.r.Body)
	if err != nil {
		return err
	}

	// Unmarshal the request body into new object of type TMSkill
	tmSkill := gormmodel.TMSkill{}
	err = json.Unmarshal(body, &tmSkill)
	if err != nil {
		return errors.MarshalingError(err)
	}

	// Validate fields of new TMSkill object
	err = c.validateTMSkillFields(&tmSkill)
	if err != nil {
		return err
	}

	// Validate that ID points to existing TMSkill in database
	tmskillSaved := gormmodel.QueryTMSKill(tmSkill.ID)
	err = c.first(&tmskillSaved)
	if err != nil {
		return err
	}

	updateMap := util.NewFilterMap("proficiency", tmSkill.Proficiency)
	err = c.updates(&tmSkill, updateMap)
	if err != nil {
		return errors.SavingError(err)
	}
	return nil
}

// Creates new TMSkill in database for POST requests to "/tmskills"
func (c *TMSkillsController) addTMSkill() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	// Unmarshal the request body into new object of type TMSkill
	tmSkill := gormmodel.TMSkill{}
	err := json.Unmarshal(body, &tmSkill)
	if err != nil {
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
	}
	// Validate fields of the TMSkill
	err = c.validateTMSkillFields(&tmSkill)
	if err != nil {
		return err
	}

	err = c.create(&tmSkill)
	if err != nil {
		return errors.SavingError(err)
	}

	// Return object JSON as response
	b, err := json.Marshal(tmSkill)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	c.Printf("Saved TMSkill: %d", tmSkill.ID)
	return nil
}

/*
validateTMSkillFields ensures that each of the following criteria are true for
the TMSkill that is passed-in:
  * the SkillID and TeamMemberID fields are both populated (not empty).
	* the SkillID and TeamMemberID fields contain the UUID of existing Skills and
	  TeamMembers in the database.
  * the Proficiency field contains a value between 0 and 5.
*/
func (c *TMSkillsController) validateTMSkillFields(tmSkill *gormmodel.TMSkill) error {
	// Validate that SkillID and TeamMemberID fields exist.
	if tmSkill.SkillID == 0 || tmSkill.TeamMemberID == 0 {
		return errors.InvalidDataModelState(fmt.Errorf(
			"A TMSkill must be a JSON object and must contain values for the %q and %q fields.",
			"skill_id", "team_member_id"))
	}

	// Validate that the IDs point to valid data.
	skill := gormmodel.QuerySkill(tmSkill.ID)
	err := c.first(&skill)
	if err != nil {
		return errors.InvalidDataModelState(fmt.Errorf(
			"the %q field of all TMSkills must contain ID of an existing Skill "+
				"in the database", "skill_id"))
	}
	teammember := gormmodel.QueryTeamMember(tmSkill.ID)
	err = c.first(&teammember)
	if err != nil {
		return errors.InvalidDataModelState(fmt.Errorf(
			"the %q field of all TMSkills must contain ID of an existing TeamMember"+
				" in the database", "team_member_id"))
	}
	// Validate that the proficiency is within the required range.
	if tmSkill.Proficiency < 0 || tmSkill.Proficiency > 5 {
		return errors.InvalidDataModelState(fmt.Errorf(
			"the %q field for a TMSkill must contain a value between 0 and 5",
			"proficiency"))
	}
	return nil
}
