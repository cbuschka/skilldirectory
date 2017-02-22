package controller

import (
	"fmt"
	"net/http"
	"skilldirectory/data"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

type BaseController struct {
	w http.ResponseWriter
	r *http.Request
	*logrus.Logger
	session    data.DataAccess
	db         *gorm.DB
	fileSystem data.FileSystem
	testSwitch bool
	errSwitch  bool
}

func (bc *BaseController) Init(w http.ResponseWriter, r *http.Request,
	session data.DataAccess, fs data.FileSystem, logger *logrus.Logger) {
	bc.w = w
	bc.r = r
	bc.Logger = logger
	bc.session = session
	bc.fileSystem = fs
}

func (bc *BaseController) InitWithGorm(w http.ResponseWriter, r *http.Request,
	session data.DataAccess, fs data.FileSystem, logger *logrus.Logger, db *gorm.DB) {
	bc.w = w
	bc.r = r
	bc.Logger = logger
	bc.session = session
	bc.fileSystem = fs
	bc.db = db
}

// GetDefaultMethods returns a string containing a ", " seperated list of the
// default HTTP methods for an endpoint.
func GetDefaultMethods() string {
	return "GET, POST, DELETE, OPTIONS"
}

// GetDefaultHeaders returns s string containing a ", " seperated list of the
// default HTTP methods for an endpoint.
func GetDefaultHeaders() string {
	return "Origin, Accept, X-Requested-With, Content-Type, " +
		"Access-Control-Request-Methods, Access-Control-Request-Headers, " +
		"Access-Control-Allow-Methods"
}

func (b *BaseController) SetTest(errSwitch bool) {
	b.testSwitch = true
	b.errSwitch = errSwitch
}

func (b BaseController) Create(object interface{}) error {
	if b.errSwitch {
		return fmt.Errorf("Error Test")
	} else if b.testSwitch {
		return nil
	}
	return b.db.Create(object).Error
}

// Delete calls gorm Delete.  Don't forget to assign the object an ID
func (b BaseController) Delete(object interface{}) error {
	if b.errSwitch {
		return fmt.Errorf("Error Test")
	} else if b.testSwitch {
		return nil
	}
	return b.db.Delete(object).Error
}
