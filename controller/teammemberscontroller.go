package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/data"
	"skilldirectory/errors"
	"skilldirectory/gormmodel"
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

func (c TeamMembersController) Options() error {
	c.w.Header().Set("Access-Control-Allow-Headers", GetDefaultHeaders())
	c.w.Header().Set("Access-Control-Allow-Methods", GetDefaultMethods())
	return nil
}

func (c *TeamMembersController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllTeamMembers()
	}

	teamMemberId, err := util.StringToID(path)
	if err != nil {
		return err
	}
	return c.getTeamMember(teamMemberId)
}

func (c *TeamMembersController) getTeamMember(id uint) error {
	teamMember := gormmodel.QueryTeamMember(id)
	err := c.first(&teamMember)
	if err != nil {
		return err
	}
	b, err := json.Marshal(teamMember)
	c.w.Write(b)
	return err
}

func (c *TeamMembersController) getAllTeamMembers() error {
	var teamMembers []gormmodel.TeamMember
	err := c.find(&teamMembers)
	if err != nil {
		return err
	}
	b, err := json.Marshal(teamMembers)
	c.w.Write(b)
	return err
}

func (c *TeamMembersController) loadTeamMember(id string) (*model.TeamMember, error) {
	teamMember := model.TeamMember{}
	err := c.session.Read("teammembers", id, data.CassandraQueryOptions{}, &teamMember)
	if err != nil {
		return nil, errors.NoSuchIDError(fmt.Errorf(
			"no TeamMeber exists with specified ID: %q", id))
	}
	return &teamMember, nil
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

	teamMember := gormmodel.TeamMember{}
	err := json.Unmarshal(body, &teamMember)
	if err != nil {
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
	}

	err = c.validatePOSTBody(&teamMember)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError type
	}

	// Save to database
	err = c.create(&teamMember)
	if err != nil {
		return errors.SavingError(err)
	}

	// Return object JSON as response
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
func (c *TeamMembersController) validatePOSTBody(teamMember *gormmodel.TeamMember) error {
	if teamMember.Name == "" || teamMember.Title == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"A Team Member must be a JSON object and must contain values for"+
				" %q and %q fields.", "name", "title"))
	}
	return nil
}
