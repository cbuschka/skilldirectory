package data

import (
	"fmt"
	"skilldirectory/util"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresConnector struct {
	db *gorm.DB
	*log.Logger
	path     string
	port     string
	keyspace string
	username string
	password string
}

type GormInterface interface {
	DB() *gorm.DB
}

func NewPostgresConnector(path, port, keyspace, username,
	password, ssl string) *PostgresConnector {
	logger := util.LogInit()
	logger.Printf("New Connector Path: %s, Port: %s, Keyspace: %s, Username: %s",
		path, port, keyspace, username)
	logger.Debug("Using Password: " + password)
	logger.Printf("New Connector Path: %s, Port: %s, Keyspace: %s, Username: %s",
		path, port, keyspace, username)
	logger.Debug("Using password: " + password)
	sslString := "sslmode=disable "
	if ssl == "true" {
		sslString = ""
	}
	postgresString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s %spassword=%s",
		path, port, username, keyspace, sslString, password)
	logger.Debugf("Postgres Connection String: %s", postgresString)

	db, err := gorm.Open("postgres", postgresString)
	if err != nil {
		fmt.Printf("Error Type: %T\n", err)
		logger.Panic(err)
	}
	return &PostgresConnector{
		db,
		logger,
		path,
		port,
		keyspace,
		username,
		password,
	}
}

func (p PostgresConnector) DB() *gorm.DB {
	return p.db
}
