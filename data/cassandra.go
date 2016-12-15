package data

import (
	"encoding/json"
	"log"

	"github.com/gocql/gocql"
)

type CassandraConnector struct {
	*gocql.Session
	path     string
	port     string
	keyspace string
}

func NewCassandraConnector(path, port, keyspace string) *CassandraConnector {
	cluster := gocql.NewCluster(path)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()

	cassConn := CassandraConnector{
		path:     path,
		port:     port,
		keyspace: keyspace,
	}
	cassConn.Session = session

	return &cassConn
}

func (c CassandraConnector) Save(key string, object interface{}) error {
	log.Print("Saving")
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return c.Query("INSERT INTO skills JSON '" + string(b) + "'").Exec()
}

func (c CassandraConnector) Read(key string, object interface{}) error {
	return nil
}

func (c CassandraConnector) Delete(key string) error {
	return nil
}

func (c CassandraConnector) ReadAll(path string, readType ReadAllInterface) ([]interface{}, error) {
	return nil, nil
}

func (c CassandraConnector) FilteredReadAll(path string, readType ReadAllInterface,
	filterFunc func(interface{}) bool) ([]interface{}, error) {
	return nil, nil
}
