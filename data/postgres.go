package data

import (
	"fmt"

	"skilldirectory/util"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresConnector struct {
	DB *gorm.DB
	*logrus.Logger
	path     string
	port     string
	keyspace string
	username string
	password string
}

func NewPostgresConnector(path, port, keyspace, username,
	password string) *PostgresConnector {
	logger := util.LogInit()
	logger.Printf("New Connector Path: %s, Port: %s, Keyspace: %s, Username: %s",
		path, port, keyspace, username)
	logger.Debug("Using password: " + password)

	DB, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		path, port, username, keyspace, password))
	if err != nil {
		logger.Panic(err)
	}
	return &PostgresConnector{
		DB,
		logger,
		path,
		port,
		keyspace,
		username,
		password,
	}
}

// // Table-Key determination is actually automatic for GORM, just print and save
// func (c PostgresConnector) Save(table, key string, object interface{}) error {
// 	c.Debugf("Saving %T (key = %s) into %s\n", object, key, table)
// 	return c.db.Save(&object).Error
// }
//
// //Reading could be just as simple if the ID is known
// func (c PostgresConnector) Read(table, id string, opts QueryOptions,
// 	object interface{}) error {
// 	c.Debugf("Reading %T (id = %s) from %s\n", object, id, table)
// 	return c.db.First(&object).Error
// }
//
// func (c PostgresConnector) ReadAll(table string, readType ReadAllInterface,
// 	targets ...[]interface{}) ([]interface{}, error) {
// 	//Just make an empty options object and pass it down
// 	return c.FilteredReadAll(table, QueryOptions{}, readType, targets...)
// }
//
// //Get all of em
// func (c PostgresConnector) FilteredReadAll(table string, opts QueryOptions,
// 	readType ReadAllInterface, targets ...[]interface{}) ([]interface{}, error) {
// 	var objects []interface{}
// 	c.Debugf("Reading all %ss from %s", readType.GetType(), table)
//
// 	db := c.db
// 	for _, filter := range opts.Filters {
// 		str := fmt.Sprintf("%s = %s", filter.key, filter.value)
// 		db = db.Where(str)
// 	}
//
// 	// I don't want to change the interface yet, so I am going to do
// 	// this next bit in two branching methods
// 	if targets != nil {
// 		//This is the preferred path, just populate the passed slice and return the
// 		//error, if there is one
// 		objects = targets[0]
// 		return nil, db.Find(&objects).Error
// 	}
//
// 	// If we're doing things the old way
//
// 	err := db.Table(table).Find(&objects).Error
// 	return objects, err
// }
//
// // If we are passed an object, deleting is cake, otherwise we can query the table
// func (c PostgresConnector) Delete(table, id string, opts QueryOptions,
// 	targets ...interface{}) error {
// 	var object interface{}
// 	if targets != nil {
// 		//We were given an object, just use it in the delete.
// 		object = targets[0]
// 	} else {
// 		//No object passed? No problem, we'll get it by querying the table.
// 		c.db.Table(table).Where("id = ?", id).First(&object)
// 	}
// 	c.Debugf("Deleting %T (id = %s) from %s\n", object, id, table)
// 	return c.db.Delete(&object).Error
// }
