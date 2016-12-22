package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Printf("removeLink() failed for the following reason:\n\t%q\n", err)
		return &errors.NoSuchIDError{
			ErrorMsg: "No Link Exists with Specified ID: " + linkID,
		}
	}

	log.Printf("Link Deleted with ID: %s", linkID)
	return nil
}

func (c *LinksController) addLink() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	link := model.Link{}
	err := json.Unmarshal(body, &link)
	if err != nil {
		return &errors.MarshalingError{
			ErrorMsg: err.Error(),
		}
	}

	err = c.validatePOSTBody(&link)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError type
	}

	if !model.IsValidLinkType(link.LinkType) {
		return &errors.InvalidLinkTypeError{
			ErrorMsg: fmt.Sprintf("Invalid Link Type: %s", link.LinkType),
		}
	}

	link.ID = util.NewID()
	err = c.session.Save("links", link.ID, link)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
	log.Printf("Saved link: %s", link.Name)
	return nil
}

/*
validatePOSTBody() accepts a model.Link pointer. It can be used to verify the
validity of the state of a Link initialized via unmarshaled JSON. Ensures that the
passed-in Link contains a key-value pair for "Name", "LinkType", "SkillID", and "URL"
fields. Returns nil error if it does, IncompletePOSTBodyError error if not.
*/
func (c *LinksController) validatePOSTBody(link *model.Link) error {
	if link.Name == "" || link.LinkType == "" ||
		link.SkillID == "" || link.URL == "" {
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: fmt.Sprintf("The JSON in a POST Request for new Link must contain values for "+
				"%q, %q, %q, and %q fields", "name", "link_type", "skill_id", "url"),
		}
	}
	return nil
}
