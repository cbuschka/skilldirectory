package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"skilldirectory/model"

	"skilldirectory/errors"

	"github.com/satori/go.uuid"
)

func loadSkill(id string) (*model.Skill, error) {
	skill := model.Skill{}
	err := skillsConnector.Read(id, &skill)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func SkillsHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Handling Skills Request: %s", r.Method)
	var err error
	var statusCode int
	switch r.Method {
	case http.MethodGet:
		err = performGet(w, r)
		statusCode = http.StatusNotFound

	case http.MethodPost:
		err = addSkill(r)
		statusCode = http.StatusBadRequest

	case http.MethodDelete:
		err = removeSkill(r)

		switch err.(type) {
		case *errors.NoSuchIDError:
			statusCode = http.StatusNotFound
		case *errors.MissingSkillIDError:
			statusCode = http.StatusBadRequest
		}
	}
	if err != nil {
		log.Printf("SkillsHandler Method: %s, Err: %v", r.Method, err)
		http.Error(w, err.Error(), statusCode)
	}
}

func addSkill(r *http.Request) error {
	// Read the body of the HTTP request into an array of bytes; ignore any errors
	body, _ := ioutil.ReadAll(r.Body)

	skill := model.Skill{}
	err := json.Unmarshal(body, &skill)
	if err != nil {
		return err
	}
	if !model.IsValidSkillType(skill.SkillType) {
		return fmt.Errorf("Invalid Skill Type: %s", skill.SkillType)
	}
	skill.Id = uuid.NewV1().String()
	err = skillsConnector.Save(skill.Id, skill)
	if err != nil {
		return err
	}
	log.Printf("New skill saved")
	return nil
}

// Removes the skill with the ID at the end of the specified request's URL.
// Returns non-nil error if the request's URL contains no ID, or if no skills
// exist with that ID.
func removeSkill(r *http.Request) error {
	// Get the ID at end of the specified request; return BadRequest400Error if request contains no ID
	skillID := checkForId(r.URL)
	if skillID == "" {
		return &errors.MissingSkillIDError{
			ErrorMsg: "No Skill ID Specified in Request URL: " + r.URL.String(),
		}
	}

	// Remove the skill with the specified ID from the skills
	// database/repository; return NoSuchID404Error if no skills have that ID
	err := skillsConnector.Delete(skillID)
	if err != nil {
		return &errors.NoSuchIDError{
			ErrorMsg: "No Skill Exists with Specified ID: " + skillID,
		}
	}

	// The skill was successfully deleted!
	log.Printf("Skill Deleted with ID: %s", skillID)
	return nil
}

func performGet(w http.ResponseWriter, r *http.Request) error {
	path := checkForId(r.URL)
	if path == "" {
		return getAllSkills(w)
	}
	return getSkill(w, path)
}

func getAllSkills(w http.ResponseWriter) error {
	skills, err := skillsConnector.ReadAll("skills/", model.Skill{})
	if err != nil {
		return err
	}
	b, err := json.Marshal(skills)
	w.Write(b)
	return err
}

func getSkill(w http.ResponseWriter, id string) error {
	skill, err := loadSkill(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(skill)
	w.Write(b)
	return err
}

// checkForID checks to see if an ID (e.g. 59317629-bcc3-11e6-9f43-6c4008bcfa84)
// has been appended to the end of the specified URL. If one has, then that ID
// will be returned. If not, then an empty string is returned ("").
func checkForId(url *url.URL) string {
	base := path.Base(url.RequestURI())
	if url.EscapedPath() != "/skills" && url.EscapedPath() != "/skills/" {
		return base
	}
	return ""
}
