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
	"strconv"
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

func (c *SkillReviewsController) getAllSkillReviews() (err error) {
	var skillReviews []model.SkillReview
	// Return all skill reviews based on skill_id in query string, if one exists
	queries := c.r.URL.Query()
	if skillID := queries.Get("skill_id"); skillID != "" {
		skillReviews, err = c.loadReviewsForSkill(skillID)
	} else { // Otherwise, return recently posted skill reviews
		// Defaults
		offsetInt := int64(0)
		limitInt := int64(10)

		// Override defaults based on query string, if one exists
		offset := queries.Get("offset")
		limit := queries.Get("limit")
		if offset != "" || limit != "" {
			var err1, err2 error
			offsetInt, err1 = strconv.ParseInt(offset, 10, 32)
			limitInt, err2 = strconv.ParseInt(limit, 10, 32)
			if err1 != nil || err2 != nil {
				return errors.InvalidQueryStringError(fmt.Errorf("Must specify a valid" +
					"integer value for offset and limit fields in query string."))
			}
		}
		skillReviews, err = c.loadRecentReviews(int(offsetInt), int(limitInt))
	}
	if err != nil {
		return err
	}

	// Respond to request with DTOs for the reviews we found:
	skillReviewDTOs := c.makeDTOs(skillReviews)
	b, err := json.Marshal(skillReviewDTOs)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)
	return nil
}

// loadReviewsForSkill returns a []model.SkillReview containing all reviews in
// the database for the skill with specified ID:
func (c *SkillReviewsController) loadReviewsForSkill(id string) ([]model.SkillReview,
	error) {
	opts := data.NewCassandraQueryOptions("skill_id", id, true)
	skillReviewsInterface, err := c.session.FilteredReadAll("skillreviews_by_skill",
		opts, model.SkillReview{})
	if err != nil {
		c.Warnf("Failed to read review from 'skillreviews_by_skill' table: %s", err)
		return nil, errors.ReadError(err)
	}
	skillReviews, err := c.convToSkillReviewSlice(skillReviewsInterface)
	if err != nil {
		return nil, err
	}
	return skillReviews, nil
}

// TODO: Would be good to offload some of this pagination work to Cassandra and
// the CQL queries themselves. But that seems to be a non-trivial task that
// Cassandra isn't really designed for.
// NOTE: The result of this function may not be what you expect if a year has
// gone by without any reviews being made. E.g. if reviews were made in 2015,
// not made in 2016, and made again in 2017; and a offset and limit are specified
// such that the function would need to get reviews from 2017 and 2015 to satify
// the limit, then only reviews from 2017 will be returned.
func (c *SkillReviewsController) loadRecentReviews(offset, limit int) ([]model.SkillReview,
	error) {
	var allReviews []model.SkillReview // Accumulates reviews from each iteration
	for i := 0; len(allReviews) < offset+limit; i++ {
		year := util.YearAsTimestamp(time.Now().Year() - i)
		opts := data.NewCassandraQueryOptions("year", year, false)
		skillReviewsInterface, err := c.session.FilteredReadAll("skillreviews_by_year",
			opts, model.SkillReview{})
		if err != nil {
			c.Warnf("Failed to read review from 'skillreviews_by_year' table: %s", err)
			return nil, errors.ReadError(err)
		}
		if len(skillReviewsInterface) == 0 {
			return allReviews[offset:], nil
		}
		skillReviews, err := c.convToSkillReviewSlice(skillReviewsInterface)
		if err != nil {
			return nil, err
		}
		allReviews = append(allReviews, skillReviews...)
	}
	return allReviews[offset : offset+limit], nil
}

// convToSkillReviewSlice converts specified []interface{} to []model.SkillReview{}
func (c *SkillReviewsController) convToSkillReviewSlice(sri []interface{}) ([]model.SkillReview,
	error) {
	skillReviewsRaw, err := json.Marshal(sri)
	if err != nil {
		c.Warnf("Encountered marshalling error: %s", err)
		return nil, errors.MarshalingError(err)
	}
	skillReviews := []model.SkillReview{}
	err = json.Unmarshal(skillReviewsRaw, &skillReviews)
	if err != nil {
		c.Warnf("Encountered marshalling error: %s", err)
		return nil, errors.MarshalingError(err)
	}
	return skillReviews, nil
}

// makeDTOs returns a slice of DTOs for the specified skill reviews:
func (c *SkillReviewsController) makeDTOs(srs []model.SkillReview) []model.SkillReviewDTO {
	skillReviewDTOs := []model.SkillReviewDTO{}
	for idx := 0; idx < len(srs); idx++ {
		skillName, err := c.getSkillName(&srs[idx])
		if err != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}
		teamMemberName, err := c.getTeamMemberName(&srs[idx])
		if err != nil {
			c.Warnf("Possible invalid id: %v", err)
			continue
		}
		skillReviewDTOs = append(skillReviewDTOs,
			srs[idx].NewSkillReviewDTO(skillName, teamMemberName))
	}
	return skillReviewDTOs
}

func (c *SkillReviewsController) getSkillReview(id string) error {
	// Get skill review from database
	skillReview, err := c.loadSkillReview(id)
	if err != nil {
		c.Warnf("Failed to read skill review: %s", err)
		return err
	}

	// Get team member and skill names associated with the review
	teamMemberName, err := c.getTeamMemberName(skillReview)
	if err != nil {
		c.Warnf("Possible invalid id: %v", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no TeamMember exists with specified ID: %q", skillReview.TeamMemberID))
	}
	skillName, err := c.getSkillName(skillReview)
	if err != nil {
		c.Warnf("Possible invalid id: %v", err)
		return errors.NoSuchIDError(fmt.Errorf(
			"no Skill exists with specified ID: %q", skillReview.SkillID))
	}

	// Send the review (and associated team member and skill names) as response
	skillReviewDTO := skillReview.NewSkillReviewDTO(skillName, teamMemberName)
	b, err := json.Marshal(skillReviewDTO)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)
	return nil
}

func (c *SkillReviewsController) loadSkillReview(id string) (*model.SkillReview,
	error) {
	skillReview := model.SkillReview{}
	err := c.session.Read("skillreviews_by_id", id, data.CassandraQueryOptions{},
		&skillReview)
	if err != nil {
		return nil,
			errors.NoSuchIDError(fmt.Errorf("no skill review exists with ID: %s", id))
	}
	return &skillReview, nil
}

func (c *SkillReviewsController) getTeamMemberName(sr *model.SkillReview) (string,
	error) {
	teamMember := model.TeamMember{}
	err := c.session.Read("teammembers", sr.TeamMemberID,
		data.CassandraQueryOptions{}, &teamMember)
	if err != nil {
		return "", err
	}
	return teamMember.Name, nil
}

func (c *SkillReviewsController) getSkillName(sr *model.SkillReview) (string,
	error) {
	skill := model.Skill{}
	err := c.session.Read("skills", sr.SkillID, data.CassandraQueryOptions{},
		&skill)
	if err != nil {
		return "", err
	}
	return skill.Name, nil
}

func (c *SkillReviewsController) removeSkillReview() error {
	// Extract required data from request body:
	body, _ := ioutil.ReadAll(c.r.Body)
	skillReview := model.SkillReview{}
	err := json.Unmarshal(body, &skillReview)
	if err != nil {
		return errors.MarshalingError(fmt.Errorf(
			"Failed to unmarshal DELETE request body: %s", err))
	}

	// Validate before deleting from database:
	err = c.validateDELETEBody(&skillReview)
	if err != nil {
		return errors.InvalidDELETEBodyError(fmt.Errorf(
			"Validation of DELETE request body failed: %s", err))
	}

	// Delete review from database:
	err = c.deleteSkillReview(skillReview.ID, skillReview.SkillID,
		skillReview.Timestamp, skillReview.Year)
	if err != nil {
		return err
	}
	log.Printf("SkillReview Deleted with ID: %s", skillReview.ID)
	return nil
}

func (c *SkillReviewsController) updateSkillReview() error {
	// Need an ID in the URL so we know what to update
	skillReviewID := util.CheckForID(c.r.URL)
	if skillReviewID == "" {
		return errors.MissingIDError(fmt.Errorf(
			"must specify a SkillReview ID in PUT request URL"))
	}

	// Load the review with ID in URL from database
	skillReview, err := c.loadSkillReview(skillReviewID)
	if err != nil {
		return err
	}

	// Update Body of skillReview with value of JSON "body" field from request's body
	bodyBytes, err := ioutil.ReadAll(c.r.Body)
	if err != nil {
		return errors.ReadError(err)
	}
	var bodyStr struct{ Body string }
	json.Unmarshal(bodyBytes, &bodyStr)
	skillReview.Body = bodyStr.Body

	// Validate before saving to database
	err = c.validatePUTBody(skillReview)
	if err != nil {
		return err
	}

	// Save updated review to database
	err = c.saveSkillReview(skillReview)
	if err != nil {
		return errors.SavingError(err)
	}
	return nil // Successfully updated review!
}

func (c *SkillReviewsController) addSkillReview() error {
	// Read body of HTTP request into model.SkillReview struct
	body, _ := ioutil.ReadAll(c.r.Body)
	skillReview := model.SkillReview{}
	err := json.Unmarshal(body, &skillReview)
	if err != nil {
		c.Warnf("Failed to unmarshal skill review: %s", err)
		return errors.MarshalingError(fmt.Errorf(
			"Failed to unmarshal POST request body: %s", err))
	}

	// Validate before saving:
	err = c.validatePOSTBody(&skillReview)
	if err != nil {
		c.Warnf("Invalid POST body: %s", err)
		return err
	}

	// Save review to database
	skillReview.ID = util.NewID()
	err = c.saveSkillReview(&skillReview)
	if err != nil {
		return errors.SavingError(err)
	}

	// Send the ID of the saved review as response
	b, _ := json.Marshal(skillReview.ID)
	c.w.Write(b)
	c.Infof("Saved SkillReview: %s", skillReview.ID)
	return nil // Successfully added review!
}

func (c *SkillReviewsController) saveSkillReview(sr *model.SkillReview) error {
	// Save model.SkillReview object to the database
	sr.Timestamp = time.Now().Format(data.TimestampFormat) // CSQL-compatible timestamp format
	sr.Year = util.YearAsTimestamp(time.Now().Year())
	err := c.session.Save("skillreviews_by_skill", sr.ID, sr)
	if err != nil {
		c.Warnf("Failed to save skill review to 'skillreviews_by_skill' table: %s", err)
		return errors.SavingError(err)
	}
	err = c.session.Save("skillreviews_by_year", sr.ID, sr)
	if err != nil {
		c.Warnf("Failed to save skill review to 'skillreviews_by_year' table: %s", err)
		return errors.SavingError(err)
	}
	err = c.session.Save("skillreviews_by_id", sr.ID, sr)
	if err != nil {
		c.Warnf("Failed to save skill review to 'skillreviews_by_id' table: %s", err)
		return errors.SavingError(err)
	}
	return nil // Successfully saved review!
}

// deleteSkillReview delete the reviews whose fields equal the specified values
// from all three skillreview tabels
func (c *SkillReviewsController) deleteSkillReview(id, skillID, timestamp,
	year string) error {
	err := c.session.Delete("skillreviews_by_id", id, data.CassandraQueryOptions{})
	if err != nil {
		c.Warnf("Failed to delete skill review from 'skillreviews_by_id' table: %s", err)
		return errors.DeleteError(err)
	}
	opts1 := data.NewCassandraQueryOptions("skill_id", skillID, true)
	opts1.AddFilter("timestamp", timestamp, false)
	err = c.session.Delete("skillreviews_by_skill", "", opts1)
	if err != nil {
		c.Warnf("Failed to delete skill review from 'skillreviews_by_skill' table: %s", err)
		return errors.DeleteError(err)
	}
	opts2 := data.NewCassandraQueryOptions("year", year, false)
	opts2.AddFilter("timestamp", timestamp, false)
	err = c.session.Delete("skillreviews_by_year", "", opts2)
	if err != nil {
		c.Warnf("Failed to delete skill review from 'skillreviews_by_year' table: %s", err)
		return errors.DeleteError(err)
	}
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

// validateDELETEBody returns error if specified SkillReview does not contain
// values for all fields required for a complete database delete operation
func (c *SkillReviewsController) validateDELETEBody(sr *model.SkillReview) error {
	if sr.ID == "" || sr.SkillID == "" || sr.Timestamp == "" || sr.Year == "" {
		return errors.InvalidDELETEBodyError(fmt.Errorf(
			"The JSON in a DELETE request to /skillreviews must contain valid values" +
				" for id, skill_id, timestamp, and year fields"))
	}
	return nil
}
