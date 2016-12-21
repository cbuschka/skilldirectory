package data

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

type CassandraConnector struct {
	*gocql.Session
	path     string
	port     string
	keyspace string
}

type CassandraQueryOptions struct {
	Filters []Filter
}

type Filter struct {
	key   string
	value string
	id    bool
}

func (f Filter) query() string {
	queryString := fmt.Sprintf(" %s", f.key)
	queryString += " = "
	if !f.id {
		queryString += "'"
	}
	queryString += f.value
	if !f.id {
		queryString += "'"
	}
	return queryString
}

func NewCassandraConnector(path, port, keyspace string) *CassandraConnector {
	log.Printf("New Connector Path: %s, Port: %s, Keyspace: %s", path, port, keyspace)
	cluster := gocql.NewCluster(path)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	cassConn := CassandraConnector{
		path:     path,
		port:     port,
		keyspace: keyspace,
	}
	cassConn.Session = session
	return &cassConn
}

// NewCassandraQueryOptions creates a new options object.
/*
key: query field name
value: query value
id: True if field is a UUID Cassandra key
*/
func NewCassandraQueryOptions(key, value string, id bool) CassandraQueryOptions {
	filter := Filter{key, value, id}
	return CassandraQueryOptions{
		Filters: []Filter{filter},
	}
}

// AddFilter adds a filter to an CassandraQueryOptions object
/*
key: query field name
value: query value
id: True if field is a UUID Cassandra key
*/
func (o *CassandraQueryOptions) AddFilter(key, value string, id bool) {
	filter := Filter{key, value, id}
	o.Filters = append(o.Filters, filter)
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
	return c.Query("DELETE FROM " + table + " WHERE id = " + key).Exec()
}

func (c CassandraConnector) ReadAll(table string, readType ReadAllInterface) ([]interface{}, error) {
	return c.FilteredReadAll(table, CassandraQueryOptions{}, readType)
}

func (c CassandraConnector) FilteredReadAll(table string, opts CassandraQueryOptions, readType ReadAllInterface) ([]interface{}, error) {
	query := "SELECT JSON * FROM " + table
	if opts.Filters != nil {
		query += " WHERE "
	}
	for _, filter := range opts.Filters {
		query += filter.query()
	}
	query += ";"
	queryBytes := []byte{}
	queryObject := readType.GetType()
	queryObjectArray := []interface{}{}
	var err error
	log.Printf("Performing query: %s", query)
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
