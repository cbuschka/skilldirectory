package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"skilldirectory/errors"
	"skilldirectory/model"
	"skilldirectory/util"
)

// SkillReviewsController handles SkillReview Requests
type SkillReviewsController struct {
	*BaseController
}

// Base implemented
func (c SkillReviewsController) Base() *BaseController {
	return c.BaseController
}

// Get implemented
func (c SkillReviewsController) Get() error {
	return c.performGet()
}

// Post implemented
func (c SkillReviewsController) Post() error {
	return c.addSkillReview()
}

// Delete implemented
func (c SkillReviewsController) Delete() error {
	return c.removeSkillReview()
}

// Put implemented
func (c SkillReviewsController) Put() error {
	return c.updateSkillReview()
}

// Options implemented
func (c SkillReviewsController) Options() error {
	c.w.Header().Set("Access-Control-Allow-Headers", GetDefaultHeaders())
	c.w.Header().Set("Access-Control-Allow-Methods", "PUT, "+GetDefaultMethods())
	return nil
}

func (c *SkillReviewsController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllSkillReviews()
	}

	skillID, err := c.pathToID(c.r.URL)
	if err != nil {
		return err
	}
	return c.getSkillReview(skillID)
}

func (c *SkillReviewsController) getAllSkillReviews() error {
	var skillReviews []model.SkillReview
	err := c.preloadAndFind(&skillReviews, "TeamMember", "Skill")
	if err != nil {
		return err
	}

	b, err := json.Marshal(skillReviews)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)
	return nil
}

func (c *SkillReviewsController) getSkillReview(id uint) error {
	skillReview := model.QuerySkillReview(id)
	err := c.preloadAndFind(&skillReview, "TeamMember", "Skill")
	if err != nil {
		return err
	}

	b, err := json.Marshal(skillReview)
	c.w.Write(b)
	return err
}

func (c *SkillReviewsController) removeSkillReview() error {
	skillID, err := c.pathToID(c.r.URL)
	if err != nil {
		return err
	}

	skillReview := model.QuerySkill(skillID)
	err = c.delete(&skillReview)
	if err != nil {
		log.Printf("removeSkillReview() failed for the following reason:"+
			"\n\t%q\n", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no SkillReview exists with specified ID: %d", skillReview.ID))
	}

	log.Printf("SkillReview Deleted with ID: %d", skillReview.ID)
	return nil
}

func (c *SkillReviewsController) updateSkillReview() error {
	skillReviewID, err := util.PathToID(c.r.URL)
	if err != nil {
		return err
	}

	skillReviewSaved := model.QuerySkillReview(skillReviewID)
	err = c.first(&skillReviewSaved)
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(c.r.Body)
	if err != nil {
		return err
	}

	type bodyStruct struct {
		Body string
	}
	var body bodyStruct
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		c.Warn(err)
	}
	var skillReviewUpdates model.SkillReview
	skillReviewUpdates.Body = body.Body

	err = c.validatePUTBody(&skillReviewUpdates)
	if err != nil {
		return err
	}
	skillReview := model.QuerySkillReview(skillReviewID)

	updateMap := util.NewFilterMap("body", skillReviewUpdates.Body)
	err = c.updates(&skillReview, updateMap)
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
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
	}

	err = c.validatePOSTBody(&skillReview)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError or errors.InvalidPOSTBodyError type
	}
	skill := model.QuerySkill(skillReview.SkillID)
	err = c.append(&skill, &skillReview, "SkillReviews")
	if err != nil {
		return errors.SavingError(err)
	}

	// Return review JSON as response
	b, err := json.Marshal(skillReview)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	log.Printf("Saved SkillReview: %d", skillReview.ID)
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
	if skillReview.SkillID == 0 ||
		skillReview.TeamMemberID == 0 ||
		skillReview.Body == "" {
		return errors.IncompletePOSTBodyError(fmt.Errorf(
			"A SkillReview must be a JSON object and must contain values for"+
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
