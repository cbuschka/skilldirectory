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
	"time"
)

type SkillReviewsController struct {
	*BaseController
}

func (c SkillReviewsController) Base() *BaseController {
	return c.BaseController
}

func (c SkillReviewsController) Get() error {
	return c.performGet()
}

func (c SkillReviewsController) Post() error {
	return c.addSkillReview()
}

func (c SkillReviewsController) Delete() error {
	return c.removeSkillReview()
}

func (c SkillReviewsController) Put() error {
	return c.updateSkillReview()
}

func (c *SkillReviewsController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllSkillReviews()
	}
	return c.getSkillReview(path)
}

func (c *SkillReviewsController) getAllSkillReviews() error {
	skillReview, err := c.session.ReadAll("skillreviews", model.SkillReview{})
	if err != nil {
		return err
	}
	b, err := json.Marshal(skillReview)
	c.w.Write(b)
	return err
}

func (c *SkillReviewsController) getSkillReview(id string) error {
	skillReview, err := c.loadSkillReview(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(skillReview)
	c.w.Write(b)
	return err
}

func (c *SkillReviewsController) loadSkillReview(id string) (*model.SkillReview, error) {
	skillReview := model.SkillReview{}
	err := c.session.Read("skillreviews", id, &skillReview)
	if err != nil {
		log.Printf("loadSkillReview() generated the following error:\n\t%q", err)
		return nil, errors.NoSuchIDError(fmt.Errorf(
			"no SkillReview exists with specified ID: %s", id))
	}
	return &skillReview, nil
}

func (c *SkillReviewsController) removeSkillReview() error {
	// Get the ID at end of the specified request
	skillReviewID := util.CheckForID(c.r.URL)
	if skillReviewID == "" {
		return errors.MissingIDError(fmt.Errorf("no SkillReview ID in request URL"))
	}

	err := c.session.Delete("skillreviews", skillReviewID, data.CassandraQueryOptions{})
	// TODO Add skillid field to opts
	if err != nil {
		log.Printf("removeSkillReview() failed for the following reason:"+
			"\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no SkillReview exists with specified ID: %s", skillReviewID))
	}

	log.Printf("SkillReview Deleted with ID: %s", skillReviewID)
	return nil
}

func (c *SkillReviewsController) updateSkillReview() error {
	skillReviewID := util.CheckForID(c.r.URL)
	if skillReviewID == "" {
		return errors.MissingIDError(fmt.Errorf(
			"must specify a SkillReview ID in PUT request URL"))
	}

	skillReview, err := c.loadSkillReview(skillReviewID)
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(c.r.Body)
	if err != nil {
		return err
	}

	var bodyStr struct{ Body string }
	json.Unmarshal(bodyBytes, &bodyStr)

	skillReview.Body = bodyStr.Body
	err = c.validatePUTBody(skillReview)
	if err != nil {
		return err
	}

	skillReview.Timestamp = time.Now().Format(data.TimestampFormat)
	err = c.session.Save("skillreviews", skillReviewID, skillReview)
	if err != nil {
		return errors.SavingError(err)
	}
	return nil
}

func (c *SkillReviewsController) addSkillReview() error {
	// Read the body of the HTTP request into an array of bytes
	body, _ := ioutil.ReadAll(c.r.Body)

	skillReview := model.SkillReview{}
	err := json.Unmarshal(body, &skillReview)
	if err != nil {
		return errors.MarshalingError(err)
	}

	err = c.validatePOSTBody(&skillReview)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError or errors.InvalidPOSTBodyError type
	}

	skillReview.Timestamp = time.Now().Format(data.TimestampFormat) // CSQL-compatible timestamp format
	skillReview.ID = util.NewID()
	err = c.session.Save("skillreviews", skillReview.ID, skillReview)
	if err != nil {
		return errors.SavingError(err)
	}
	log.Printf("Saved SkillReview: %s", skillReview.ID)
	return nil
}

/*
validatePOSTBody() accepts a model.SkillReview pointer. It can be used to verify the
validity of the state of a SkillReview initialized via unmarshaled JSON. Ensures that
the passed-in SkillReview contains a key-value pair for "SkillID", "TeamMemberID",
"Body", and "Date" fields. Returns nil error if it does, IncompletePOSTBodyError
error if not.
*/
func (c *SkillReviewsController) validatePOSTBody(skillReview *model.SkillReview) error {
	if skillReview.SkillID == "" || skillReview.TeamMemberID == "" ||
		skillReview.Body == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"The JSON in a POST Request for new SkillReview must contain values for"+
				" %q, %q, and %q fields.", "skill_id", "team_member_id", "body"))
	}
	return nil
}

/*
validatePUTBody() accepts a model.SkillReview pointer. It can be used to verify the
validity of the state of a SkillReview updated via a PUT request. Ensures that the
passed-in SkillReview's "Body" field is not empty.
*/
func (c *SkillReviewsController) validatePUTBody(skillReview *model.SkillReview) error {
	if skillReview.Body == "" {
		return errors.InvalidPUTBodyError(fmt.Errorf(
			"The JSON in a PUT request for new SkillReview must contain a value "+
				"for the %q field", "body"))
	}
	return nil
}
