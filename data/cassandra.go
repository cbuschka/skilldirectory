package data

import (
	"encoding/json"

	"github.com/gocql/gocql"
)

type CassandraConnector struct {
	*gocql.Session
	path     string
	port     string
	table    string
	keyspace string
}

func NewCassandraConnector(path, port, table, keyspace string) *CassandraConnector {
	cluster := gocql.NewCluster(path)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()

	cassConn := CassandraConnector{
		path:     path,
		port:     port,
		keyspace: keyspace,
		table:    table,
	}
	cassConn.Session = session

	return &cassConn
}

func (c CassandraConnector) Save(key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return c.Query("INSERT INTO %s ", c.table, b).Exec()
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
