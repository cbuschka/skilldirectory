package controller

import (
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"skilldirectory/util"
	"skilldirectory/model"
	"skilldirectory/errors"
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
	return fmt.Errorf("PUT requests nor currently supported.")
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
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No SkillReview Exists with Specified ID: " + id,
		}
	}
	return &skillReview, nil
}

func (c *SkillReviewsController) removeSkillReview() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillReviewID := util.CheckForID(c.r.URL)
	if skillReviewID == "" {
		return &errors.MissingIDError{
			ErrorMsg: "No SkillReview ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete("skillreviews", skillReviewID, "team_member_id")
	if err != nil {
		log.Printf("removeSkillReview() failed for the following reason:\n\t%q\n", err)
		return &errors.NoSuchIDError{
			ErrorMsg: "No SkillReview Exists with Specified ID: " + skillReviewID,
		}
	}

	log.Printf("SkillReview Deleted with ID: %s", skillReviewID)
	return nil
}

func (c *SkillReviewsController) addSkillReview() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	skillReview := model.SkillReview{}
	err := json.Unmarshal(body, &skillReview)
	if err != nil {
		return &errors.MarshalingError{
			ErrorMsg: err.Error(),
		}
	}

	err = c.validatePOSTBody(&skillReview)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError or errors.InvalidPOSTBodyError type
	}

	skillReview.ID = util.NewID()
	err = c.session.Save("skillreviews", skillReview.ID, skillReview)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
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
	   skillReview.Body    == "" || skillReview.Date         == "" {
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: fmt.Sprintf("The JSON in a POST Request for new SkillReview must contain values for" +
				" %q, %q, %q, and %q fields.", "skill_id", "team_member_id", "body", "date"),
		}
	}
	return nil
}
