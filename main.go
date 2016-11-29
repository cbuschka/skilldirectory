package main

import (
	"net/http"
	"skilldir/router"
)

// var skillsConnector = data.NewAccessor(data.NewFileWriter("skills/"))

func main() {
	// addSkillsToView()
	router := router.StartRouter()
	http.ListenAndServe(":8080", router)
}

// func addSkillsToView() {
// 	skillsConnector.Save("Golang", model.NewSkill("Golang", "language"))
// }
