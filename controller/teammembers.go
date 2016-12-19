package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"skilldirectory/errors"
	"skilldirectory/model"

	"github.com/satori/go.uuid"
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
	path := checkForId(c.r.URL)
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
	b, err := json.Marshal(teamMember)
	c.w.Write(b)
	return err
}

func (c *TeamMembersController) loadTeamMember(id string) (*model.TeamMember, error) {
	teamMember := model.TeamMember{}
	err := c.session.Read("teammembers", id, &teamMember)
	if err != nil {
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No Team Member Exists with Specified ID: " + id,
		}
	}
	return &teamMember, nil
}

func (c *TeamMembersController) removeTeamMember() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	teamMemberID := checkForId(c.r.URL)
	if teamMemberID == "" {
		return &errors.MissingIDError{
			ErrorMsg: "No Team Member ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete("teammembers", teamMemberID)
	if err != nil {
		return &errors.NoSuchIDError{
			ErrorMsg: "No Team Member Exists with Specified ID: " + teamMemberID,
		}
	}

	log.Printf("Team Member Deleted with ID: %s", teamMemberID)
	return nil
}

func (c *TeamMembersController) addTeamMember() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	teamMember := model.TeamMember{}
	err := json.Unmarshal(body, &teamMember)
	if err != nil {
		return &errors.MarshalingError{
			ErrorMsg: err.Error(),
		}
	}

	err = c.validatePOSTBody(&teamMember)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError type
	}

	teamMember.ID = uuid.NewV1().String()
	err = c.session.Save("teammembers", teamMember.ID, teamMember)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
	log.Printf("Saved Team Member: %s", teamMember.Name)
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
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: "POST Request for new Team Member must contain values for " +
				"\"Name\" and \"Title\" fields.",
		}
	}
	return nil
}
