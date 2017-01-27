package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/data"
	"skilldirectory/errors"
	"skilldirectory/model"
	"skilldirectory/util"
)

type SkillIconsController struct {
	*BaseController
}

func (c SkillIconsController) Base() *BaseController {
	return c.BaseController
}

func (c SkillIconsController) Get() error {
	return c.performGet()
}

func (c SkillIconsController) Post() error {
	return c.addSkillIcon()
}

func (c SkillIconsController) Delete() error {
	return c.removeSkillIcon()
}

func (c SkillIconsController) Put() error {
	return c.addSkillIcon()
}

func (c SkillIconsController) performGet() error {
	path := util.CheckForID(c.r.URL)
	if path == "" {
		return c.getAllSkillIcons()
	}
	return c.getSkillIcon(path)
}

func (c *SkillIconsController) getAllSkillIcons() error {
	skillIcons, err := c.session.ReadAll("skillIcons", model.SkillIcon{})
	if err != nil {
		return err
	}

	b, err := json.Marshal(skillIcons)
	c.w.Write(b)
	return err
}

func (c *SkillIconsController) getSkillIcon(skillID string) error {
	skillIcon := model.SkillIcon{}
	err := c.session.Read("skillicons", "",
		data.NewCassandraQueryOptions("skill_id", skillID, true), &skillIcon)
	if err != nil {
		return errors.NoSuchIDError(
			fmt.Errorf("no skill icon exists for skill with ID: %s", skillID))
	}

	b, err := json.Marshal(skillIcon)
	c.w.Write(b)
	return err
}

func (c *SkillIconsController) removeSkillIcon() error {
	// Get ID at end of request; return error if request contains no ID
	skillID := util.CheckForID(c.r.URL)
	if skillID == "" {
		return errors.MissingIDError(fmt.Errorf("no skill ID specified in request URL"))
	}

	// Attempt to delete image resource from S3
	err := c.fileSystem.Delete("dev/" + skillID)
	if err != nil {
		c.Warn(err)
		return err
	}

	// Attempt to delete record from database
	err = c.session.Delete("skillicons", "",
		data.NewCassandraQueryOptions("skill_id", skillID, true))
	if err != nil {
		c.Warnf("Failed to delete skill icon from database.")
		return errors.NoSuchIDError(fmt.Errorf(
			"no skill icon exists with specified ID: %s", skillID))
	}

	c.Printf("SkillIcon Deleted with ID: %s", skillID)
	return nil
}

// Creates new SkillIcon in database for POST requests to "/skillicons"
func (c *SkillIconsController) addSkillIcon() error {
	// Extract icon data from HTTP request
	iconFile, _, err := c.r.FormFile("icon")
	if err != nil {
		c.Warn("error getting icon form file: " + err.Error())
		return errors.ReadError(fmt.Errorf("Failed to parse icon field: %s", err))
	}
	defer iconFile.Close()

	// Unmarshal the request body into new object of type SkillIcon
	skillIcon := model.SkillIcon{
		SkillID: c.r.FormValue("skill_id"),
	}

	// Capture data for later use before it is consumed by util.ValidateIcon
	iconFileBytes, _ := ioutil.ReadAll(iconFile)

	// Validity and error checking
	dataCopy := make([]byte, len(iconFileBytes))
	copy(dataCopy, iconFileBytes)
	_, err = util.ValidateIcon(bytes.NewReader(dataCopy))
	if err != nil {
		c.Warn("Invalid image data: ", err)
		return errors.InvalidPOSTBodyError(err)
	}
	err = validateSkillID(skillIcon.SkillID, c.session)
	if err != nil {
		c.Warn("ID does not exist: ", err.Error())
		return errors.InvalidPOSTBodyError(fmt.Errorf(
			"The %q field must contain ID of existing Skill in database", "skill_id"))
	}

	// Upload image to S3 cloud
	url, err := c.fileSystem.Write("dev/"+skillIcon.SkillID,
		bytes.NewReader(iconFileBytes))
	if err != nil {
		return fmt.Errorf("failed to save icon: %s", err)
	}

	// Store its URL in Cassandra table
	skillIcon.URL = url
	err = c.session.Save("skillicons", skillIcon.SkillID, skillIcon)
	if err != nil {
		return errors.SavingError(err)
	}

	c.Printf("Saved icon: %s", skillIcon.URL)
	return nil
}
