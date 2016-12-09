package controller

import (
	"net/http"
	"skilldirectory/data"
)

type BaseController struct {
	w       http.ResponseWriter
	r       *http.Request
	session data.DataAccess
}

func NewController(w http.ResponseWriter, r http.Request, d data.DataAccess) BaseController {
	return BaseController{w, &r, d}
}
