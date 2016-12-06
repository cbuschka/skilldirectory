package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

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
	}
	if err != nil {
		log.Printf("SkillsHandler Method: %s, Err: %v", r.Method, err)
		http.Error(w, err.Error(), statusCode)
	}
}

func addSkill(r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	skill := model.Skill{}
	err = json.Unmarshal(body, &skill)
	if err != nil {
		log.Println(err)
		return err
	}
	if !model.IsValidSkillType(skill.SkillType) {
		return fmt.Errorf("Invalid Skill Type: %s", skill.SkillType)
	}
	skill.Id = uuid.NewV1().String()
	err = skillsConnector.Save(skill.Id, skill)
	if err != nil {
		log.Printf("Save skill: %s error: %s", skill.Name, err)
	}
	log.Printf("New skill saved")
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
	skills := []model.Skill{}
	filepath.Walk("skills/", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			skill, err1 := loadSkill(f.Name())
			if err1 != nil {
				log.Println(err1)
				return err1
			}
			skills = append(skills, *skill)
		}
		return nil
	})
	b, err := json.Marshal(skills)
	if err != nil {
		log.Printf("Marshal skills error: %s", err.Error())
		return err
	}
	w.Write(b)
	return nil
}

func getSkill(w http.ResponseWriter, id string) error {
	skill, err := loadSkill(id)
	if err != nil {
		return err
	}
	b, err := json.Marshal(skill)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func checkForId(url *url.URL) string {
	_, path := path.Split(url.RequestURI())
	return path
}
