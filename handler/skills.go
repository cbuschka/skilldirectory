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
	"skilldir/model"
)

func loadSkill(title string) (*model.Skill, error) {
	skill := model.Skill{}
	err := skillsConnector.Read(title, &skill)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func SkillsHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Handling Skills Request: %s", r.Method)
	switch r.Method {
	case http.MethodGet:
		err := performGet(w, r)
		if err != nil {
			log.Printf("getSkills: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	case http.MethodPost:
		err := addSkill(r)
		if err != nil {
			log.Printf("addSkill: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
	err = skillsConnector.Save(skill.Name, skill)
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
	fmt.Println("Return one skill")
	return nil
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

func getSkill(w http.ResponseWriter, path string) {

}

// Returns true if there is an id
func checkForId(url *url.URL) string {
	_, path := path.Split(url.RequestURI())
	return path
}
