package data

import (
  "skilldirectory/util"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "github.com/Sirupsen/logrus"
)

type PostgresConnector struct {
  *gorm.DB
  *logrus.Logger
  path      string
  port      string
  keyspace  string
  username  string
  password  string
}

func NewPostgresConnector (path, port, keyspace, username,
   password string) *PostgresConnector {
   logger := util.LogInit()
   logger.Printf("New Connector Path: %s, Port: %s, Keyspace: %s, Username: %s",
    path, port, keyspace, username)
   logger.Debug("Using password: " + password)

   db, err := gorm.Open("postgres",
     "host=%s:%s user=%s dbname=%s sslmode=disable password=%s",
     path, port, username, keyspace, password)
   if err != nil {
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
