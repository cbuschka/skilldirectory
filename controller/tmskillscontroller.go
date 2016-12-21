package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	tmSkill, err := c.session.ReadAll("tmskills", model.TMSkill{})
	if err != nil {
		return err
	}
	b, err := json.Marshal(tmSkill)
	c.w.Write(b)
	return err
}

func (c *TMSkillsController) getTMSkill(id string) error {
	tmSkill, err := c.loadTMSkill(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(tmSkill)
	c.w.Write(b)
	return err
}

func (c *TMSkillsController) loadTMSkill(id string) (*model.TMSkill, error) {
	tmSkill := model.TMSkill{}
	err := c.session.Read("tmskills", id, &tmSkill)
	if err != nil {
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No TMSkill Exists with Specified ID: " + id,
		}
	}
	return &tmSkill, nil
}

func (c *TMSkillsController) removeTMSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	tmSkillID := util.CheckForID(c.r.URL)
	if tmSkillID == "" {
		return &errors.MissingIDError{
			ErrorMsg: "No TMSkill ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete("tmskills", tmSkillID)
	if err != nil {
		return &errors.NoSuchIDError{
			ErrorMsg: "No TMSkill Exists with Specified ID: " + tmSkillID,
		}
	}

	log.Printf("TMSkill Deleted with ID: %s", tmSkillID)
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
		return err // Will be of errors.IncompletePOSTBodyError or errors.InvalidPOSTBodyError type
	}

	tmSkill.ID = util.NewID()
	err = c.session.Save("tmskills", tmSkill.ID, tmSkill)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
	log.Printf("Saved TMSkill: %s", tmSkill.ID)
	return nil
}

/*
validatePOSTBody() accepts a model.TMSkill pointer. It can be used to verify the
validity of the state of a TMSkill initialized via unmarshaled JSON. Ensures that

	* the passed-in TMSkill contains a key-value pair for "SkillID" and for "TeamMemberID"
	fields. Returns nil error if it does, IncompletePOSTBodyError error if not.

	* the value specified for the "Proficiency" field is between 0 and 5. Returns nil error
	if it is, InvalidPOSTBodyError if not.
*/
func (c *TMSkillsController) validatePOSTBody(tmSkill *model.TMSkill) error {
	if tmSkill.SkillID == "" || tmSkill.TeamMemberID == "" {
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: "The JSON in a POST Request for new TMSkill must contain values for " +
				"\"skill_id\" and \"team_member_id\" fields.",
		}
	}
	if tmSkill.Proficiency < 0 || tmSkill.Proficiency > 5 {
		return &errors.InvalidPOSTBodyError{
			ErrorMsg: "The JSON in a POST Request for new TMSkill must contain \"proficiency\" " +
				"field with a value between 0 and 5.",
		}
	}
	return nil
}
