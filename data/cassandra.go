package data

import (
	"encoding/json"
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
	logger := util.LogInit()
	logger.Printf("New Connector Path: %s, Port: %s, Keyspace: %s", path, port, keyspace)
	cluster := gocql.NewCluster(path)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Panic(err)
	}
	cassConn := CassandraConnector{
		path:     path,
		port:     port,
		keyspace: keyspace,
	}
	cassConn.Session = session
	cassConn.Logger = logger
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

/*
Delete will run a DELETE query on the database, thereby removing the specified record.
Delete requires a table to delete from, and an id to identify the record. If the table
uses a compound PRIMARY KEY, then it is necessary to specify the column names of all
additional columns used in the PRIMARY KEY besides "id" (it is assumed that "id" is used
in the PRIMARY KEY).

Note that Delete will still be able to execute the DELETE query if "id" is specified as a primary_key_column.
*/
func (c CassandraConnector) Delete(table, id string, primary_key_cols ...string) error {
	query := "DELETE FROM " + table + " WHERE id = " + id // Base query

	// If the PRIMARY KEY is compound (uses more columns than "id"), then append those onto the base query
	for _, col := range primary_key_cols {
		if col == "id" { // Skip the "id" column; it's already in the base query
			continue
		}
		// Get value for this column
		colQuery := "SELECT " + col + " FROM " + table + " WHERE id = " + id
		m := make(map[string]interface{})
		err := c.Query(colQuery).Consistency(gocql.One).MapScan(m)
		if err != nil {
			return err
		}

		// value can be of different types. Based on type, append it onto the
		// base query string.
		var colVal string
		switch v := m[col].(type) {
		case int:
			colVal = string(v)
		case string:
			colVal = "'" + v + "'"
		case gocql.UUID:
			colVal = v.String()
		default:
			return fmt.Errorf("Unrecognized value type, %T, for column, %q", m[col], col)
		}
		query += (" AND " + col + " = " + colVal)
	}

	c.Printf("Running the following DELETE query:\n\t%q\n", query)
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
