package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		b, err := getSkills()
		if err != nil {
			log.Printf("getSkills: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write(b)
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

func getSkills() ([]byte, error) {
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
		return nil, err
	}
	return b, nil
}
