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
		statusCode = http.StatusNotFound
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
	// Get the ID at end of the specified request; return error if request contains no ID
	skillID := checkForId(r.URL)
	if skillID == "" {
		return fmt.Errorf("No Skill ID Specified in Request URL: %s", r.URL)
	}

	// Remove the skill with the specified ID from the skills
	// database/repository; return error if no skills have that ID
	err := skillsConnector.Delete(skillID)
	if err != nil {
		return fmt.Errorf("No Skill Exists with Specified ID: %s", skillID)
	}

	// The skill was successfully deleted!
	log.Printf("Skill Deleted with ID: %s", skillID)
	return nil
}

func performGet(w http.ResponseWriter, r *http.Request) error {
	path := checkForId(r.URL)
	if path == "" {
		filter, err := extractSkillFilter(r.URL)
		if err != nil {
			return err
		}

		if filter == "" {
			return getAllSkills(w)
		} else {
			return getAllSkillsFiltered(w, filter)
		}
	}
	return getSkill(w, path)
}

func getAllSkills(w http.ResponseWriter) error {
	// Read all skills from the database/repository.
	// Store the read's results in an array of model.Skill{}s.
	skills, err := skillsConnector.ReadAll("skills/", model.Skill{})
	if err != nil {
		return err
	}

	// Marshal the array of skills into JSON format, and send
	// it in a response through the passed-in ResponseWriter.
	b, err := json.Marshal(skills)
	w.Write(b)
	return err
}

func getAllSkillsFiltered(w http.ResponseWriter, filter string) error {
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
	filteredSkills, err := skillsConnector.FilteredReadAll("skills/", model.Skill{}, filterer)
	if err != nil {
		return err
	}

	// Encode the slice into JSON format and send it in a response via the passed-in ResponseWriter
	b, err := json.Marshal(filteredSkills)
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

func extractSkillFilter(url *url.URL) (string, error) {
	// Extract the URL path's base and make sure it is "/skills".
	// Return error if it's not, because it doesn't make sense to filter
	// by skill type if the base isn't "/skills".
	base := path.Base(url.Path)
	if base != "skills" {
		return "", fmt.Errorf("URL path base must be \"skills\" to filter by skill type")
	}

	// Extract the query string from the URL as a key, value map.
	// Then search the map for a "skilltype" filter. If the query
	// string contains this filter, return the value. If not, return
	// an empty string ("").
	query := url.Query()
	for key, val := range query {
		if key == "skilltype" {
			return val[0], nil
		}
	}
	return "", nil
}
