package handler

import (
	"encoding/json"
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
		skills := getSkills()
		b, err := json.Marshal(skills)
		if err != nil {
			log.Printf("Marshal skills error: %s", err.Error())
			return
		}
		w.Write(b)

	}
}

func getSkills() []model.Skill {
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
	return skills
}
