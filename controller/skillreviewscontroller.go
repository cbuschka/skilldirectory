package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"skilldirectory/errors"
	"skilldirectory/gormmodel"
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

func (c SkillReviewsController) Options() error {
	c.w.Header().Set("Access-Control-Allow-Headers", GetDefaultHeaders())
	c.w.Header().Set("Access-Control-Allow-Methods", "PUT, "+GetDefaultMethods())
	return nil
}

//
func (c *SkillReviewsController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllSkillReviews()
	}

	skillID, err := util.StringToID(path)
	if err != nil {
		return err
	}
	return c.getSkillReview(skillID)
}

func (c *SkillReviewsController) getAllSkillReviews() error {
	var skillReviews []gormmodel.SkillReview
	err := c.find(&skillReviews)
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

// func (c *SkillReviewsController) getReviewsForSkill(skill_id string) error {
// 	opts := data.NewCassandraQueryOptions("skill_id", skill_id, false)
// 	skillReviewsInterface, err := c.session.FilteredReadAll("skillreviews", opts, model.SkillReview{})
// 	if err != nil {
// 		return err
// 	}
// 	skillReviews, err := convertToReviewsStruct(skillReviewsInterface)
// 	if err != nil {
// 		return errors.MarshalingError(err)
// 	}
//
// 	skillReviewDTOs := c.convertReviewsToDTOs(skillReviews)
//
// 	b, err := json.Marshal(skillReviewDTOs)
// 	if err != nil {
// 		return errors.MarshalingError(err)
// 	}
// 	c.w.Write(b)
// 	return err
// }
//
func (c *SkillReviewsController) getSkillReview(id uint) error {
	var skillReview gormmodel.SkillReview
	err := c.first(&skillReview)
	if err != nil {
		return err
	}

	// teamMemberName, err := c.getTeamMemberName(skillReview)
	// if err != nil {
	// 	c.Warnf("Possible invalid id: %v", err)
	// 	return errors.NoSuchIDError(fmt.Errorf(
	// 		"no TeamMember exists with specified ID: %q", skillReview.TeamMemberID))
	// }
	//
	// skillName, err := c.getSkillName(skillReview)
	// if err != nil {
	// 	c.Warnf("Possible invalid id: %v", err)
	// 	return errors.NoSuchIDError(fmt.Errorf(
	// 		"no Skill exists with specified ID: %q", skillReview.SkillID))
	// }

	// skillReviewDTO := skillReview.NewSkillReviewDTO(skillName, teamMemberName)
	b, err := json.Marshal(skillReview)
	c.w.Write(b)
	return err
}

// func (c *SkillReviewsController) loadSkillReview(id string) (*gormmodel.SkillReview,
// 	error) {
// 	skillReview := gormmodel.SkillReview{}
// 	err := c.first(&skillReview)
// 	if err != nil {
// 		log.Printf("loadSkillReview() generated the following error:\n\t%q", err)
// 		return nil, errors.NoSuchIDError(fmt.Errorf(
// 			"no SkillReview exists with specified ID: %d", id))
// 	}
// 	return &skillReview, nil
// }

// func (c *SkillReviewsController) getTeamMemberName(sr *model.SkillReview) (string,
// 	error) {
// 	teamMember := model.TeamMember{}
// 	err := c.session.Read("teammembers", sr.TeamMemberID,
// 		data.CassandraQueryOptions{}, &teamMember)
// 	if err != nil {
// 		return "", err
// 	}
// 	return teamMember.Name, nil
// }
//
// func (c *SkillReviewsController) getSkillName(sr *model.SkillReview) (string,
// 	error) {
// 	skill := model.Skill{}
// 	err := c.session.Read("skills", sr.SkillID, data.CassandraQueryOptions{},
// 		&skill)
// 	if err != nil {
// 		return "", err
// 	}
// 	return skill.Name, nil
// }

func (c *SkillReviewsController) removeSkillReview() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return errors.MissingIDError(fmt.Errorf("Missing required id for DELETE call"))
	}

	skillID, err := util.StringToID(path)
	if err != nil {
		return err
	}

	skillReview := gormmodel.QuerySkill(skillID)
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
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return errors.MissingIDError(fmt.Errorf(
			"must specify a SkillReview ID in PUT request URL"))
	}

	skillReviewId, err := util.StringToID(path)
	if err != nil {
		return err
	}

	skillReviewSaved := gormmodel.QuerySkillReview(skillReviewId)
	err = c.first(&skillReviewSaved)
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(c.r.Body)
	if err != nil {
		return err
	}

	var bodyStr struct {
		Body     string
		Positive bool
	}
	json.Unmarshal(bodyBytes, &bodyStr)

	skillReview := gormmodel.QuerySkillReview(skillReviewId)
	skillReview.Body = bodyStr.Body
	skillReview.Positive = bodyStr.Positive
	err = c.validatePUTBody(&skillReview)
	if err != nil {
		return err
	}

	updateMap := make(map[string]interface{})
	updateMap["body"] = bodyStr.Body
	updateMap["positive"] = bodyStr.Positive
	err = c.updates(&skillReview, updateMap)
	if err != nil {
		return errors.SavingError(err)
	}
	return nil
}

func (c *SkillReviewsController) addSkillReview() error {
	// Read the body of the HTTP request into an array of bytes
	body, _ := ioutil.ReadAll(c.r.Body)
	skillReview := gormmodel.SkillReview{}
	fmt.Println(skillReview)
	err := json.Unmarshal(body, &skillReview)
	if err != nil {
		c.Warn("Marshaling Error: ", errors.MarshalingError(err))
	}

	err = c.validatePOSTBody(&skillReview)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError or errors.InvalidPOSTBodyError type
	}
	skill := gormmodel.QuerySkill(skillReview.SkillID)
	err = c.append(&skill, skillReview, "SkillReviews")
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
func (c *SkillReviewsController) validatePOSTBody(skillReview *gormmodel.SkillReview) error {
	fmt.Println(skillReview)
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
func (c *SkillReviewsController) validatePUTBody(skillReview *gormmodel.SkillReview) error {
	if skillReview.Body == "" {
		return errors.InvalidPUTBodyError(fmt.Errorf(
			"The JSON in a PUT request for new SkillReview must contain a value "+
				"for the %q field", "body"))
	}
	return nil
}
