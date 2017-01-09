package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/data"
	"skilldirectory/errors"
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

func (c LinksController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllLinks()
	}
	return c.getLink(path)
}

func (c *LinksController) getAllLinks() error {
	var links []interface{}
	var err error
	filter := c.r.URL.Query().Get("linktype")
	var opts data.CassandraQueryOptions

	// Add approved query filters here
	if filter != "" {
		opts = data.NewCassandraQueryOptions("linktype", filter, false)
	}
	links, err = c.session.FilteredReadAll("links", opts, model.Link{})

	if err != nil {
		return err
	}

	b, err := json.Marshal(links)
	c.w.Write(b)
	return err
}

func (c *LinksController) getLink(id string) error {
	link, err := c.loadLink(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(link)
	c.w.Write(b)
	return err
}

func (c *LinksController) loadLink(id string) (*model.Link, error) {
	link := model.Link{}
	err := c.session.Read("links", id, &link)
	if err != nil {
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No Link Exists with Specified ID: " + id,
		}
	}
	return &link, nil
}

func (c *LinksController) removeLink() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	linkID := util.CheckForID(c.r.URL)
	if linkID == "" {
		return &errors.MissingIDError{
			ErrorMsg: "No Link ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete("links", linkID, "skill_id")
	if err != nil {
		c.Printf("removeLink() failed for the following reason:\n\t%q\n", err)
		return &errors.NoSuchIDError{
			ErrorMsg: "No Link Exists with Specified ID: " + linkID,
		}
	}

	c.Printf("Link Deleted with ID: %s", linkID)
	return nil
}

// Creates new Link in database for POST requests to "/links"
func (c *LinksController) addLink() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	// Unmarshal the request body into new object of type Link
	link := model.Link{}
	err := json.Unmarshal(body, &link)
	if err != nil {
		return &errors.MarshalingError{
			ErrorMsg: err.Error(),
		}
	}

	// Validate fields of the Link
	err = c.validateLinkFields(&link)
	if err != nil {
		return err
	}

	// Save the Link to database under a new UUID
	link.ID = util.NewID()
	err = c.session.Save("links", link.ID, link)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
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
func (c *LinksController) validateLinkFields(link *model.Link) error {
	// Validate that SkillID field exists
	if link.SkillID == "" || link.LinkType == "" ||
		link.Name == "" || link.URL == "" {
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: fmt.Sprintf("the JSON in a POST Request for new Link must "+
				"contain values for %q, %q, %q, and %q fields",
				"name", "link_type", "skill_id", "url"),
		}
	}

	// Validate that SkillID points to valid data
	err := c.session.Read("skills", link.SkillID, &model.Skill{})
	if err != nil {
		return &errors.InvalidDataModelState{
			ErrorMsg: fmt.Sprintf("the %q field of all Links must contain ID of "+
				"an existing skill in the database", "skill_id"),
		}
	}

	// Validate the the LinkType field is valid
	if !model.IsValidLinkType(link.LinkType) {
		return &errors.InvalidLinkTypeError{
			ErrorMsg: fmt.Sprintf("Invalid Link Type: %q", link.LinkType),
		}
	}
	return nil
}
