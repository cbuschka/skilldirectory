package controller

import (
	"encoding/json"
	"fmt"
	"skilldirectory/data"
	"skilldirectory/errors"
	"skilldirectory/model"
	"skilldirectory/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"mime/multipart"
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
	err := c.session.Read("skillicons", skillID, &skillIcon)
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
	err := c.deleteOnAWS(skillID)
	if err != nil {
		c.Warnf("Failed to delete skill icon from AWS S3.")
		return errors.NoSuchIDError(fmt.Errorf(
			"no skill icon exists with specified ID: %s", skillID))
	}

	// Attempt to delete record from database
	err = c.session.Delete("skillicons", skillID, data.CassandraQueryOptions{})
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
	c.r.ParseMultipartForm(1500000) // using 1.5 MB buffer
	iconFile, _, err := c.r.FormFile("icon")
	if err != nil {
		c.Warn("error getting icon form file: " + err.Error())
		return errors.ReadError(fmt.Errorf("Failed to parse icon field"))
	}
	defer iconFile.Close()

	// Unmarshal the request body into new object of type SkillIcon
	skillIcon := model.SkillIcon{
		SkillID: c.r.FormValue("skill_id"),
	}

	// Validity and error checking
	_, err = util.ValidateIcon(iconFile)
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
	err = c.saveOnAWS(iconFile, skillIcon.SkillID)
	if err != nil {
		return fmt.Errorf("failed to save icon: %s", err)
	}

	// Store its URL in Cassandra table
	skillIcon.URL = "https://s3.amazonaws.com/skilldirectory/andrew/" +
		skillIcon.SkillID
	err = c.session.Save("skillicons", skillIcon.SkillID, skillIcon)
	if err != nil {
		return errors.SavingError(err)
	}

	c.Printf("Saved icon: %s", skillIcon.URL)
	return nil
}

// Stores icon in the '/skilldirectory' AWS bucket
func (c *SkillIconsController) saveOnAWS(icon multipart.File, name string) error {
	// Establish connection to AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return fmt.Errorf("failed to connect to AWS S3 instance")
	}
	svc := s3.New(sess)

	// Setup object to save in AWS
	params := &s3.PutObjectInput{
		Bucket: aws.String("skilldirectory"), // Required
		Key:    aws.String("andrew/" + name), // Required
		Body:   icon,
		Metadata: map[string]*string{
			"Andrew": aws.String("Dillon"),
		},
	}

	// Try to save the icon to AWS
	_, err = svc.PutObject(params)
	if err != nil {
		return fmt.Errorf("Failed to save icon to AWS S3 instance")
	}
	return nil
}

// Deletes icon with name in '/skilldirectory' AWS bucket
func (c *SkillIconsController) deleteOnAWS(name string) error {
	// Establish connection to AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return fmt.Errorf("failed to connect to AWS S3 instance")
	}
	svc := s3.New(sess)

	// Setup object to delete from AWS
	params := &s3.DeleteObjectInput{
		Bucket: aws.String("skilldirectory"), // Required
		Key:    aws.String("andrew/" + name), // Required
	}

	// Try to delete the icon from AWS
	_, err = svc.DeleteObject(params)
	if err != nil {
		return fmt.Errorf("Failed to delete icon from AWS S3 instance")
	}
	return nil

}
