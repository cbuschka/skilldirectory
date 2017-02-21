package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/data"
	"skilldirectory/errors"
	"skilldirectory/model"
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
	return c.getTMSkill(path)
}

func (c *TMSkillsController) getAllTMSkills() error {
	tmSkillsInterface, err := c.session.ReadAll("tmskills", model.TMSkill{})
	if err != nil {
		return err
	}

	tmSkills, err := convertToStruct(tmSkillsInterface)
	if err != nil {
		return err
	}

	tmSkillDTOs := c.convertTMSkillsToDTOs(tmSkills)

	b, err := json.Marshal(tmSkillDTOs)
	c.w.Write(b)
	return err
}

func convertToStruct(tmSkillsInterface []interface{}) ([]model.TMSkill, error) {
	tmSkillsRaw, err := json.Marshal(tmSkillsInterface)
	if err != nil {
		return nil, errors.MarshalingError(err)
	}

	tmSkills := []model.TMSkill{}
	err = json.Unmarshal(tmSkillsRaw, &tmSkills)
	if err != nil {
		return nil, errors.MarshalingError(err)
	}
	return tmSkills, nil
}

func (c *TMSkillsController) convertTMSkillsToDTOs(tmSkills []model.TMSkill) []model.TMSkillDTO {
	tmSkillDTOs := []model.TMSkillDTO{}
	for idx := 0; idx < len(tmSkills); idx++ {
		skillName, err := c.getSkillName(&tmSkills[idx])
		if err != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}

		teamMemberName, err2 := c.getTeamMemberName(&tmSkills[idx])
		if err2 != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}
		tmSkillDTOs = append(tmSkillDTOs,
			tmSkills[idx].NewTMSkillDTO(skillName, teamMemberName))
	}
	return tmSkillDTOs
}

func (c *TMSkillsController) getTMSkill(id string) error {
	tmSkill, err := c.loadTMSkill(id)
	if err != nil {
		return err
	}

	teamMemberName, err := c.getTeamMemberName(tmSkill)
	if err != nil {
		c.Warnf("Possible invalid id: %v", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no TeamMember exists with specified ID: %q", tmSkill.TeamMemberID))
	}

	skillName, err := c.getSkillName(tmSkill)
	if err != nil {
		c.Warnf("Possible invalid id: %v", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no Skill exists with specified ID: %q", tmSkill.SkillID))
	}

	tmSkillDTO := tmSkill.NewTMSkillDTO(skillName, teamMemberName)
	b, err := json.Marshal(tmSkillDTO)
	c.w.Write(b)
	return err
}

func (c *TMSkillsController) loadTMSkill(id string) (*model.TMSkill, error) {
	tmSkill := model.TMSkill{}
	err := c.session.Read("tmskills", id, data.CassandraQueryOptions{}, &tmSkill)
	if err != nil {
		c.Warnf("loadTMSkill() generated the following error: %v", err)
		return nil, errors.NoSuchIDError(fmt.Errorf(
			"No TMSkill Exists with Specified ID: %s ", id))
	}
	return &tmSkill, nil
}

func (c *TMSkillsController) getTeamMemberName(tmSkill *model.TMSkill) (string, error) {
	teamMember := model.TeamMember{}
	err := c.session.Read("teammembers", tmSkill.TeamMemberID,
		data.CassandraQueryOptions{}, &teamMember)
	if err != nil {
		return "", err
	}
	return teamMember.Name, nil
}

func (c *TMSkillsController) getSkillName(tmSkill *model.TMSkill) (string, error) {
	skill := model.Skill{}
	err := c.session.Read("skills", tmSkill.SkillID, data.CassandraQueryOptions{},
		&skill)
	if err != nil {
		return "", err
	}
	return skill.Name, nil
}

func (c *TMSkillsController) removeTMSkill() error {
	// Get the ID at end of the request; return error if request contains no ID
	tmSkillID := util.CheckForID(c.r.URL)
	if tmSkillID == "" {
		return errors.MissingIDError(fmt.Errorf("no TMSkill ID in request URL"))
	}

	err := c.session.Delete("tmskills", tmSkillID, data.NewCassandraQueryOptions("team_member_id", "", true))

	if err != nil {
		c.Printf("removeTMSkill() failed for the following reason:\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no TMSkill exists with specified ID: %q", tmSkillID))
	}

	c.Printf("TMSkill Deleted with ID: %s", tmSkillID)
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
	tmSkill := model.TMSkill{}
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
	tmSkill.ID = tmSkillID
	err = c.validateTMSkillID(&tmSkill)
	if err != nil {
		return err
	}

	// Write the updated TMSkill and return
	err = c.session.Save("tmskills", tmSkill.ID, tmSkill)
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
	tmSkill := model.TMSkill{}
	err := json.Unmarshal(body, &tmSkill)
	if err != nil {
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
	}
	// Validate fields of the TMSkill
	err = c.validateTMSkillFields(&tmSkill)
	if err != nil {
		return err
	}

	// Save the TMSkill to database under a new UUID
	tmSkill.ID = util.NewID()
	err = c.session.Save("tmskills", tmSkill.ID, tmSkill)
	if err != nil {
		return errors.SavingError(err)
	}

	// Return object JSON as response
	b, err := json.Marshal(tmSkill)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	c.Printf("Saved TMSkill: %s", tmSkill.ID)
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
func (c *TMSkillsController) validateTMSkillFields(tmSkill *model.TMSkill) error {
	// Validate that SkillID and TeamMemberID fields exist.
	if tmSkill.SkillID == "" || tmSkill.TeamMemberID == "" {
		return errors.InvalidDataModelState(fmt.Errorf(
			"A TMSkill must be a JSON object and must contain values for the %q and %q fields.",
			"skill_id", "team_member_id"))
	}

	// Validate that the IDs point to valid data.
	err := c.session.Read("skills", tmSkill.SkillID, data.CassandraQueryOptions{},
		&model.Skill{})
	if err != nil {
		return errors.InvalidDataModelState(fmt.Errorf(
			"the %q field of all TMSkills must contain ID of an existing Skill "+
				"in the database", "skill_id"))
	}
	err = c.session.Read("teammembers", tmSkill.TeamMemberID,
		data.CassandraQueryOptions{}, &model.TeamMember{})
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

/*
validateTMSkillID ensures that the ID of the specified TMSkill is the ID of an
existing TMSkill entry in the database.
*/
func (c *TMSkillsController) validateTMSkillID(tmSkill *model.TMSkill) error {
	// Validate that the TMSkill's ID exists in the database
	err := c.session.Read("tmskills", tmSkill.ID, data.CassandraQueryOptions{},
		&model.TMSkill{})
	if err != nil {
		return errors.NoSuchIDError(fmt.Errorf(
			"the following ID is not valid: %s", tmSkill.ID))
	}
	return nil
}
