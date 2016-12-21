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

func (bc *BaseController) Init(w http.ResponseWriter, r *http.Request, session data.DataAccess) {
	bc.w = w
	bc.r = r
	bc.session = session
}
