package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"skilldirectory/model"

	"skilldirectory/errors"

	"fmt"

	uuid "github.com/satori/go.uuid"
)

type SkillsController struct {
	*BaseController
}

func (c SkillsController) Base() *BaseController {
	return c.BaseController
}
func (c SkillsController) Get() error {
	return c.performGet()
}
func (c SkillsController) Post() error {
	return c.addSkill()
}

func (c SkillsController) Delete() error {
	return c.removeSkill()
}

func (c SkillsController) Put() error {
	return nil
}

func (c SkillsController) performGet() error {
	path := checkForId(c.r.URL)
	if path == "" {
		filter := c.r.URL.Query().Get("skilltype")
		if filter == "" {
			return c.getAllSkills()
		} else {
			return c.getAllSkillsFiltered(filter)
		}
	}
	return c.getSkill(path)
}

func (c *SkillsController) getAllSkills() error {
	skills, err := c.session.ReadAll("skills/", model.Skill{})
	if err != nil {
		return err
	}
	b, err := json.Marshal(skills)
	c.w.Write(b)
	return err
}

func (c *SkillsController) getAllSkillsFiltered(filter string) error {
	// Only try to apply the specified filter if it is either a valid Skill Type, or else
	// is a wildcard filter ("").
	if !model.IsValidSkillType(filter) && filter != "" {
		return fmt.Errorf("The skilltype filter, \"%s\", is not valid", filter)
	}

	// This function is used as the filter for the call to skillsConnectory.FilteredReadAll() below.
	// It compares the SkillType field of the skills read from the database/repository to the passed-in
	// filter string. Only those skills whose SkillType matches the filter string pass through.
	filterer := func(object interface{}) bool {
		// Each object that is passed in is of type map[string]interface{}, so must cast to that.
		// Then, objmap is a mapping of Skill type fields to their values.
		// For example, fmt.Println(object), might display:
		// 	map[Id:9dbdbca3-be38-11e6-bdb2-6c4008bcfa84 Name:Java SkillType:database]
		objmap := object.(map[string]interface{})
		if objmap["SkillType"] == filter {
			return true
		}
		return false
	}

	// Get a slice containing all skills from the skills database/repository that pass through the filter function.
	filteredSkills, err := c.session.FilteredReadAll("skills/", model.Skill{}, filterer)
	if err != nil {
		return err
	}

	// Encode the slice into JSON format and send it in a response via the passed-in ResponseWriter
	b, err := json.Marshal(filteredSkills)
	c.w.Write(b)
	return err
}

func (c *SkillsController) getSkill(id string) error {
	skill, err := c.loadSkill(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(skill)
	c.w.Write(b)
	return err
}

func (c *SkillsController) loadSkill(id string) (*model.Skill, error) {
	skill := model.Skill{}
	err := c.session.Read(id, &skill)
	if err != nil {
		return nil, &errors.NoSuchIDError{
			ErrorMsg: "No Skill Exists with Specified ID: " + id,
		}
	}
	return &skill, nil
}

func (c *SkillsController) removeSkill() error {
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID := checkForId(c.r.URL)
	if skillID == "" {
		return &errors.MissingSkillIDError{
			ErrorMsg: "No Skill ID Specified in Request URL: " + c.r.URL.String(),
		}
	}

	err := c.session.Delete(skillID)
	if err != nil {
		return &errors.NoSuchIDError{
			ErrorMsg: "No Skill Exists with Specified ID: " + skillID,
		}
	}

	log.Printf("Skill Deleted with ID: %s", skillID)
	return nil
}

func (c *SkillsController) addSkill() error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(c.r.Body)

	err := validatePOSTBody(body)
	if err != nil {
		return err // Will be of errors.IncompletePOSTBodyError type
	}

	skill := model.Skill{}
	err = json.Unmarshal(body, &skill)
	if err != nil {
		return &errors.MarshalingError{
			ErrorMsg: err.Error(),
		}
	}
	if !model.IsValidSkillType(skill.SkillType) {
		return &errors.InvalidSkillTypeError{
			ErrorMsg: fmt.Sprint("Invalid Skill Type: %s", skill.SkillType),
		}
	}
	skill.ID = uuid.NewV1().String()
	err = c.session.Save(skill.ID, skill)
	if err != nil {
		return &errors.SavingError{
			ErrorMsg: err.Error(),
		}
	}
	log.Printf("Saved skill: %s", skill.Name)
	return nil
}

// validatePOSTBody() accepts a byte array of a JSON request body. Ensures that the
// passed-in byte[] request contains a key-value pair for "Name" and for "Skilltype"
// fields. Returns nil error if it does, IncompletePOSTBodyError error if not.
func validatePOSTBody(body []byte) error {
	object := model.Skill{}
	json.Unmarshal(body, &object)
	if object.Name == "" || object.SkillType == "" {
		return &errors.IncompletePOSTBodyError{
			ErrorMsg: "POST Request for new Skill must contain values for " +
				"\"Name\" and \"SkillType\" fields.",
		}
	}
	return nil
}