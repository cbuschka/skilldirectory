package controller

import (
	"fmt"
	"net/http"
	"skilldirectory/data"
	"skilldirectory/gormmodel"

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

func (bc *BaseController) SetTest(errSwitch bool) {
	bc.testSwitch = true
	bc.errSwitch = errSwitch
}

func (bc BaseController) create(object gormmodel.GormInterface) error {
	if bc.errSwitch {
		return fmt.Errorf("Error Test")
	} else if bc.testSwitch {
		return nil
	}
	return bc.db.Create(object).Error
}

// Delete calls gorm Delete.  Don't forget to assign the object an ID
func (bc BaseController) delete(object gormmodel.GormInterface) error {
	if object.GetID() == 0 {
		return fmt.Errorf("Can't Delete Nil Object")
	} else if bc.errSwitch {
		return fmt.Errorf("Error Test")
	} else if bc.testSwitch {
		return nil
	}
	return bc.db.Delete(object).Error
}

func (bc BaseController) first(object gormmodel.GormInterface) error {
	if bc.errSwitch {
		return fmt.Errorf("Error Test")
	} else if bc.testSwitch {
		return nil
	}
	return bc.db.First(object, object.GetID()).Error
}

func (bc BaseController) find(object interface{}) error {
	if bc.errSwitch {
		return fmt.Errorf("Error Test")
	} else if bc.testSwitch {
		return nil
	}
	return bc.db.Where("deleted_at IS NULL").Find(object).Error
}

func (bc BaseController) preloadAndFind(object interface{}, preload ...string) error {
	if bc.errSwitch {
		return fmt.Errorf("Error Test")
	} else if bc.testSwitch {
		return nil
	}
	db := *bc.db
	for _, p := range preload {
		db = *db.Preload(p)
	}
	return db.Find(object).Error
}

func (bc BaseController) updates(object gormmodel.GormInterface, updateMap map[string]interface{}) error {
	if bc.errSwitch {
		return fmt.Errorf("Error Test")
	} else if bc.testSwitch {
		return nil
	}

	return bc.db.Model(object).Updates(updateMap).Error
}

/*
append takes a reference to a parentObject (&Skill), a childAppend value (SkillReview) and that
association sting ("SkillReviews")
*/
func (bc BaseController) append(parentObject, childAppend gormmodel.GormInterface, association string) error {
	if bc.errSwitch {
		return fmt.Errorf("Error Test")
	} else if bc.testSwitch {
		return nil
	}
	return bc.db.Model(parentObject).Association(association).Append(childAppend).Error
}
