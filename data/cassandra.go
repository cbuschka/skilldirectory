package data

import (
	"encoding/json"
	"fmt"

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

func (c CassandraConnector) Save(table, key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return c.Query("INSERT INTO " + table + " JSON '" + string(b) + "'").Exec()
}

func (c CassandraConnector) Read(table, key string, object interface{}) error {
	query := "SELECT JSON * FROM " + table + " WHERE id = " + key
	byteQ := []byte{}
	err := c.Query(query).Consistency(gocql.One).Scan(&byteQ)
	if err != nil {
		return err
	}
	return json.Unmarshal(byteQ, &object)
}

func (c CassandraConnector) Delete(table, key string) error {
	return nil
}

func (c CassandraConnector) ReadAll(table, path string, readType ReadAllInterface) ([]interface{}, error) {
	query := "SELECT JSON * FROM " + table
	queryBytes := []byte{}
	queryObject := readType.GetType()
	queryObjectArray := []interface{}{}
	var err error
	iter := c.Query(query).Iter()
	for iter.Scan(&queryBytes) {
		err = json.Unmarshal(queryBytes, &queryObject)
		if err != nil {
			return nil, fmt.Errorf("Read all unmarshal err: %s", err)
		}
		queryObjectArray = append(queryObjectArray, queryObject)
	}
	return queryObjectArray, nil
}

func (c CassandraConnector) FilteredReadAll(path string, readType ReadAllInterface,
	filterFunc func(interface{}) bool) ([]interface{}, error) {
	return nil, nil
}
