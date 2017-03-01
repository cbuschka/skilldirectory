package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/errors"
	"skilldirectory/gormmodel"
	"skilldirectory/model"
	"skilldirectory/util"
)

type LinksController struct {
	*BaseController
}

func (c LinksController) Base() *BaseController {
	return c.BaseController
}

func (c LinksController) Get() error {
	return c.performGet()
}

func (c LinksController) Post() error {
	return c.addLink()
}

func (c LinksController) Delete() error {
	return c.removeLink()
}

func (c LinksController) Put() error {
	return fmt.Errorf("PUT requests not currently supported.")
}

func (c LinksController) Options() error {
	c.w.Header().Set("Access-Control-Allow-Headers", GetDefaultHeaders())
	c.w.Header().Set("Access-Control-Allow-Methods", GetDefaultMethods())
	return nil
}

func (c LinksController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllLinks()
	}
	linkID, err := util.StringToID(path)
	if err != nil {
		return err
	}
	return c.getLink(linkID)
}

func (c *LinksController) getAllLinks() error {
	var links []gormmodel.Link
	var err error
	filter := c.r.URL.Query().Get("linktype")

	// Add approved query filters here
	if filter != "" {
		filterMap := util.NewFilterMap("link_type", filter)
		err = c.findWhere(&links, filterMap)
	} else {
		err = c.find(&links)
	}

	if err != nil {
		return err
	}

	b, err := json.Marshal(links)
	c.w.Write(b)
	return err
}

func (c *LinksController) getLink(id uint) error {
	link, err := c.loadLink(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(link)
	c.w.Write(b)
	return err
}

func (c *LinksController) loadLink(id uint) (*gormmodel.Link, error) {
	link := gormmodel.QueryLink(id)
	err := c.first(&link)
	if err != nil {
		return nil, errors.NoSuchIDError(
			fmt.Errorf("no Link exists with specified ID: %d", id))
	}
	return &link, nil
}

func (c *LinksController) removeLink() error {
	// Get ID at end of request; return error if request contains no ID
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return errors.MissingIDError(fmt.Errorf("no Link ID specified in request URL"))
	}

	linkID, err := util.StringToID(path)
	if err != nil {
		return err
	}
	link := gormmodel.QueryLink(linkID)
	err = c.delete(&link)

	if err != nil {
		c.Printf("removeLink() failed for the following reason:\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no Link exists with specified ID: %d", linkID))
	}

	c.Printf("Link Deleted with ID: %d", linkID)
	return nil
}

// Creates new Link in database for POST requests to "/links"
func (c *LinksController) addLink() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	// Unmarshal the request body into new object of type Link
	link := gormmodel.Link{}
	err := json.Unmarshal(body, &link)
	if err != nil {
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
	}
	// Validate fields of the Link
	err = c.validateLinkFields(&link)
	if err != nil {
		return err
	}

	err = c.create(&link)
	if err != nil {
		return errors.SavingError(err)
	}

	b, err := json.Marshal(link)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	c.Printf("Saved link: %s", link.Name)
	return nil
}

/*
validateLinkFields ensures that each of the following criteria are true for the
Link that is passed-in:
  * the SkillID, LinkType, Name, and URL fields are populated (not empty).
	* the SkillID field contains the UUID of an existing Skill in the database.
	* the LinkType field contains valid link type (see model.IsValidLinkType)
*/
func (c *LinksController) validateLinkFields(link *gormmodel.Link) error {
	// Validate that SkillID field exists
	if link.SkillID == 0 || link.LinkType == "" ||
		link.Name == "" || link.URL == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"A Link must be a JSON object and must contain values for "+
				"%q, %q, %q, and %q fields", "name", "link_type", "skill_id", "url"))
	}

	// Validate that SkillID points to valid data
	skill := gormmodel.QuerySkill(link.SkillID)
	err := c.first(&skill)
	if err != nil {
		return errors.InvalidDataModelState(fmt.Errorf(
			"the %q field of all Links must contain ID of an existing skill in "+
				"the database", "skill_id"))
	}

	// Validate the the LinkType field is valid
	if !model.IsValidLinkType(link.LinkType) {
		return errors.InvalidLinkTypeError(fmt.Errorf(
			"Invalid Link Type: %q", link.LinkType))
	}
	return nil
}
