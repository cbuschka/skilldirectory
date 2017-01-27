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

type TeamMembersController struct {
	*BaseController
}

func (c TeamMembersController) Base() *BaseController {
	return c.BaseController
}

func (c TeamMembersController) Get() error {
	return c.performGet()
}

func (c TeamMembersController) Post() error {
	return c.addTeamMember()
}

func (c TeamMembersController) Delete() error {
	return c.removeTeamMember()
}

func (c TeamMembersController) Put() error {
	return fmt.Errorf("PUT requests nor currently supported.")
}

func (c *TeamMembersController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllTeamMembers()
	}
	return c.getTeamMember(path)
}

func (c *TeamMembersController) getAllTeamMembers() error {
	teamMembers, err := c.session.ReadAll("teammembers", model.TeamMember{})
	if err != nil {
		return err
	}
	b, err := json.Marshal(teamMembers)
	c.w.Write(b)
	return err
}

func (c *TeamMembersController) getTeamMember(id string) error {
	teamMember, err := c.loadTeamMember(id)
	if err != nil {
		return err
	}

	tmSkillDTOs, err := c.getAllTMSkills(teamMember)
	if err != nil {
		return err
	}

	teamMemberDTO := teamMember.NewTeamMemberDTO(tmSkillDTOs)
	b, err := json.Marshal(teamMemberDTO)
	c.w.Write(b)
	return err
}

func (c *TeamMembersController) loadTeamMember(id string) (*model.TeamMember, error) {
	teamMember := model.TeamMember{}
	err := c.session.Read("teammembers", id, &teamMember)
	if err != nil {
		return nil, errors.NoSuchIDError(fmt.Errorf(
			"no TeamMeber exists with specified ID: %q", id))
	}
	return &teamMember, nil
}

// Get a []model.TMSkillDTO for all TMSkills in the database associated with the
// specified TeamMember
func (c *TeamMembersController) getAllTMSkills(teamMember *model.TeamMember) (
	[]model.TMSkillDTO, error) {
	// Get all TMSkills that reference the passed-in TeamMember
	options := data.NewCassandraQueryOptions(
		"team_member_id", teamMember.ID, false)
	tmSkillsInterface, err := c.session.FilteredReadAll("tmskills",
		options, model.TMSkill{})
	if err != nil {
		return nil, err
	}

	// Convert byte representation to map[string]interface{}
	tmSkillsRaw, err := json.Marshal(tmSkillsInterface)
	if err != nil {
		return nil, err
	}

	// Convert map[string]interface{} to []model.TMSkill{}
	tmSkills := []model.TMSkill{}
	err = json.Unmarshal(tmSkillsRaw, &tmSkills)
	if err != nil {
		return nil, err
	}

	// "Convert" []model.TMSkill{} to []model.TMSkillDTO{} so that the
	// names of each TMSkill's TeamMember and Skill IDs are returned.
	tmSkillsDTO := []model.TMSkillDTO{}
	for _, tmSkill := range tmSkills {
		// Get name of TeamMember, skip if encounter an error
		teamMemberName, err := c.getTeamMemberName(&tmSkill)
		// If encounter an error, log and continue to next TMSkill
		if err != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}

		// Get name of Skill, skip if encounter an error
		skillName, err := c.getSkillName(&tmSkill)
		if err != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}

		// Append new TMSkillDTO to return object w/ the names
		tmSkillsDTO = append(tmSkillsDTO,
			tmSkill.NewTMSkillDTO(skillName, teamMemberName))
	}

	return tmSkillsDTO, nil
}

func (c *TeamMembersController) getTeamMemberName(tmSkill *model.TMSkill) (string, error) {
	teamMember := model.TeamMember{}
	err := c.session.Read("teammembers", tmSkill.TeamMemberID, &teamMember)
	if err != nil {
		return "", err
	}
	return teamMember.Name, nil
}

func (c *TeamMembersController) getSkillName(tmSkill *model.TMSkill) (string, error) {
	skill := model.Skill{}
	err := c.session.Read("skills", tmSkill.SkillID, &skill)
	if err != nil {
		return "", err
	}
	return skill.Name, nil
}

func (c *TeamMembersController) removeTeamMember() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	teamMemberID := util.CheckForID(c.r.URL)
	if teamMemberID == "" {
		return errors.MissingIDError(fmt.Errorf("no TeamMember ID in request URL"))
	}

	err := c.session.Delete("teammembers", teamMemberID, data.CassandraQueryOptions{})
	if err != nil {
		c.Printf("removeTeamMember() failed for the following reason:\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"No Team Member Exists with Specified ID: %q", teamMemberID))
	}

	c.Printf("Team Member Deleted with ID: %s", teamMemberID)
	return nil
}

func (c *TeamMembersController) addTeamMember() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	teamMember := model.TeamMember{}
	err := json.Unmarshal(body, &teamMember)
	if err != nil {
		return errors.MarshalingError(err)
	}

	err = c.validatePOSTBody(&teamMember)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError type
	}

	teamMember.ID = util.NewID()
	err = c.session.Save("teammembers", teamMember.ID, teamMember)

	if err != nil {
		return errors.SavingError(err)
	}

	b, err := json.Marshal(teamMember)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	c.Infof("Saved Team Member: %s", teamMember.Name)
	return nil
}

/*
validatePOSTBody() accepts a model.TeamMember pointer. It can be used to verify the
validity of the state of a TeamMember initialized via unmarshaled JSON. Ensures that the
passed-in TeamMember contains a key-value pair for "Name" and for "Title"
fields. Returns nil error if it does, IncompletePOSTBodyError error if not.
*/
func (c *TeamMembersController) validatePOSTBody(teamMember *model.TeamMember) error {
	if teamMember.Name == "" || teamMember.Title == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"The JSON in a POST Request for new Team Member must contain values for"+
				" %q and %q fields.", "name", "title"))
	}
	return nil
}
