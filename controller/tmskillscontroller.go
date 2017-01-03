package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/errors"
	"skilldirectory/model"
	util "skilldirectory/util"
)

type TMSkillsController struct {
	*BaseController
}

func (c TMSkillsController) Base() *BaseController {
	return c.BaseController
}

func (c TMSkillsController) Get() error {
	return c.performGet()
}

func (c TMSkillsController) Post() error {
	return c.addTMSkill()
}

func (c TMSkillsController) Delete() error {
	return c.removeTMSkill()
}

func (c TMSkillsController) Put() error {
	return fmt.Errorf("PUT requests nor currently supported.")
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

	tmSkillsRaw, err := json.Marshal(tmSkillsInterface)
	if err != nil {
		return &errors.MarshalingError{err.Error()}
	}

	tmSkills := []model.TMSkill{}
	err = json.Unmarshal(tmSkillsRaw, &tmSkills)
	if err != nil {
		return &errors.MarshalingError{err.Error()}
	}

	tmSkillDTOs := []model.TMSkillDTO{}
	for idx := 0; idx < len(tmSkills); idx++ {
		skillName, err := c.getSkillName(&tmSkills[idx])
		if err != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}

		teamMemberName, err := c.getTeamMemberName(&tmSkills[idx])
		if err != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}
		tmSkillDTOs = append(tmSkillDTOs,
			tmSkills[idx].NewTMSkillDTO(skillName, teamMemberName))
	}

	b, err := json.Marshal(tmSkillDTOs)
	c.w.Write(b)
	return err
}

func (c *TMSkillsController) getTMSkill(id string) error {
	tmSkill, err := c.loadTMSkill(id)
	if err != nil {
		return err
	}

	teamMemberName, err := c.getTeamMemberName(tmSkill)
	if err != nil {
		c.Warnf("Possible invalid id: %v", err)
		return &errors.NoSuchIDError{
			ErrorMsg: fmt.Sprintf("no TeamMember exists with "+
				"specified ID: %q", tmSkill.TeamMemberID),
		}
	}

	skillName, err := c.getSkillName(tmSkill)
	if err != nil {
		c.Warnf("Possible invalid id: %v", err)
		return &errors.NoSuchIDError{
			ErrorMsg: fmt.Sprintf("no Skill exists with "+
				"specified ID: %q", tmSkill.SkillID),
		}
	}

	tmSkillDTO := tmSkill.NewTMSkillDTO(skillName, teamMemberName)
	b, err := json.Marshal(tmSkillDTO)
	c.w.Write(b)
	return err
}

func (c *TMSkillsController) loadTMSkill(id string) (*model.TMSkill, error) {
	tmSkill := model.TMSkill{}
	err := c.session.Read("tmskills", id, &tmSkill)
	if err != nil {
		c.Warnf("loadTMSkill() generated the following error: %v", err)
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No TMSkill Exists with Specified ID: " + id,
		}
	}
	return &tmSkill, nil
}

func (c *TMSkillsController) getTeamMemberName(tmSkill *model.TMSkill) (string, error) {
	teamMember := model.TeamMember{}
	err := c.session.Read("teammembers", tmSkill.TeamMemberID, &teamMember)
	if err != nil {
		return "", err
	}
	return teamMember.Name, nil
}

func (c *TMSkillsController) getSkillName(tmSkill *model.TMSkill) (string, error) {
	skill := model.Skill{}
	err := c.session.Read("skills", tmSkill.SkillID, &skill)
	if err != nil {
		return "", err
	}
	return skill.Name, nil
}

func (c *TMSkillsController) removeTMSkill() error {
	// Get the ID at end of the request; return error if request contains no ID
	tmSkillID := util.CheckForID(c.r.URL)
	if tmSkillID == "" {
		return &errors.MissingIDError{
			ErrorMsg: "No TMSkill ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete("tmskills", tmSkillID, "team_member_id")
	if err != nil {
		c.Printf("removeTMSkill() failed for the following reason:\n\t%q\n", err)
		return &errors.NoSuchIDError{
			ErrorMsg: "No TMSkill Exists with Specified ID: " + tmSkillID,
		}
	}

	c.Printf("TMSkill Deleted with ID: %s", tmSkillID)
	return nil
}

func (c *TMSkillsController) addTMSkill() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	tmSkill := model.TMSkill{}
	err := json.Unmarshal(body, &tmSkill)
	if err != nil {
		return &errors.MarshalingError{
			ErrorMsg: err.Error(),
		}
	}

	err = c.validatePOSTBody(&tmSkill)
	if err != nil {
		return err // Will be IncompletePOSTBodyError or InvalidPOSTBodyError type
	}

	tmSkill.ID = util.NewID()
	err = c.session.Save("tmskills", tmSkill.ID, tmSkill)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
	c.Printf("Saved TMSkill: %s", tmSkill.ID)
	return nil
}

/*
validatePOSTBody() accepts a model.TMSkill pointer. It can be used to verify the
validity of the state of a TMSkill initialized via unmarshaled JSON. Ensures that

	* the tmSkill contains a key-value pair for "SkillID" and "TeamMemberID"
	fields. Returns nil error if it does, IncompletePOSTBodyError error if not.

	* the value specified for the "Proficiency" field is between 0 and 5. Returns
	nil error if it is, InvalidPOSTBodyError if not.
*/
func (c *TMSkillsController) validatePOSTBody(tmSkill *model.TMSkill) error {
	// Validate existence of SkillID and TeamMemberID fields
	if tmSkill.SkillID == "" || tmSkill.TeamMemberID == "" {
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: fmt.Sprintf("The JSON in POST Request for new TMSkill must "+
				"contain values for %q and %q fields",
				"skill_id", "team_member_id"),
		}
	}

	// Validate existence of Skill with specified ID
	err := c.session.Read("skills", tmSkill.SkillID, model.Skill{})
	if err != nil {
		return &errors.NoSuchIDError{
			ErrorMsg: "No skill Exists with ID: " + tmSkill.SkillID,
		}
	}

	// Validate existence of Team Member with specified ID
	err = c.session.Read("teammembers", tmSkill.TeamMemberID, model.TeamMember{})
	if err != nil {
		return &errors.NoSuchIDError{
			ErrorMsg: "No Team Member Exists with ID: " + tmSkill.TeamMemberID,
		}
	}

	// Validate the Proficiency field value is within required range
	if tmSkill.Proficiency < 0 || tmSkill.Proficiency > 5 {
		return &errors.InvalidPOSTBodyError{
			ErrorMsg: fmt.Sprintf("The JSON in POST Request for new TMSkill must "+
				"contain %q field with value between 0 and 5",
				"proficiency"),
		}
	}

	return nil
}
