package data

/*
DataAccess represents an interface for accessing and managing a data storage
system (e.g. a database). Implementations of DataAccess should ensure that the
objects passed to them are properly serialized under, and can be retrieved by
providing, the specified key string.
*/
type DataAccess interface {
	Save(table, key string, object interface{}) error
	Read(table, key string, object interface{}) error
	ReadAll(table string, readType ReadAllInterface) ([]interface{}, error)
	FilteredReadAll(table string, opts CassandraQueryOptions,
		readType ReadAllInterface) ([]interface{}, error)
	Delete(table, id string, opts CassandraQueryOptions) error
}

/*
ReadAllInterface implementations must contain a GetType() function that returns
the type of the implementor.
*/
type ReadAllInterface interface {
	GetType() interface{}
}
