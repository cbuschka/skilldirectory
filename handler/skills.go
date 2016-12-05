package handler

import (
	"encoding/json"
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
	var err error
	switch r.Method {
	case http.MethodGet:
		err = performGet(w, r)

	case http.MethodPost:
		err = addSkill(r)
	}
	if err != nil {
		log.Printf("SkillsHandler Method: %s, Err: %v", r.Method, err)
		http.Error(w, err.Error(), http.StatusNotFound)
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
