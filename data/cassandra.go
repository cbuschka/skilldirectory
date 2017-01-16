package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"skilldirectory/util"

	"github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
)

// TimestampFormat is a CSQL-compatible timestamp format
const TimestampFormat = "2006-01-02 15:04:05-0700"

type CassandraConnector struct {
	*gocql.Session
	*logrus.Logger
	path     string
	port     string
	keyspace string
	username string
	password string
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

func NewCassandraConnector(path, port, keyspace, username,
	password string) *CassandraConnector {
	logger := util.LogInit()
	logger.Printf("New Connector Path: %s, Port: %s, Keyspace: %s, Username: %s",
		path, port, keyspace, username)
	logger.Debug("Using Password: " + password)
	cluster := gocql.NewCluster(path)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Panic(err)
	}
	cassConn := CassandraConnector{
		path:     path,
		port:     port,
		keyspace: keyspace,
		username: username,
		password: password,
	}
	cassConn.Session = session
	cassConn.Logger = logger
	return &cassConn
}

/*
NewCassandraQueryOptions creates a new options object.
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

/*
AddFilter adds a filter to an CassandraQueryOptions object
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

/*
Delete will run a DELETE query on the database, thereby removing the specified record.
Delete requires a table to delete from, and an id to identify the record. If the table
uses a compound PRIMARY KEY, then it is necessary to specify the column names of all
additional columns used in the PRIMARY KEY besides "id" (it is assumed that "id" is used
in the PRIMARY KEY).

Note that Delete will still be able to execute the DELETE query if "id" is specified as a primary_key_column.
*/
func (c CassandraConnector) Delete(table, id string, opts CassandraQueryOptions) error {
	query := "DELETE FROM " + table + " WHERE " // Base query
	c.Infof("Deleting:%s,%s,%v", table, id, opts)
	if id == "" && len(opts.Filters) == 0 {
		return errors.New("Attempting to delete with no id")
	}
	if id != "" {
		query += "id = " + id
	}
	for _, filter := range opts.Filters {
		if filter.key == "id" {
			continue
		}

		query += filter.query()
	}
	query += ";"
	c.Infof("Running the following DELETE query:\n\t%q\n", query)
	return c.Query(query).Exec()
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
	c.Debugf("Performing query: %s", query)
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
