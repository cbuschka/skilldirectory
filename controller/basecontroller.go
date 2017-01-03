package controller

import (
	"net/http"
	"skilldirectory/data"
	"skilldirectory/util"

	"github.com/Sirupsen/logrus"
)

type BaseController struct {
	w http.ResponseWriter
	r *http.Request
	*logrus.Logger
	session data.DataAccess
}

func (bc *BaseController) Init(w http.ResponseWriter, r *http.Request, session data.DataAccess) {
	bc.w = w
	bc.r = r
	bc.Logger = util.LogInit()
	bc.session = session
}
