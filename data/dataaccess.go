package data

/*
DataAccess represents an interface for accessing and managing a data storage
system (e.g. a database). Implementations of DataAccess should ensure that the
objects passed to them are properly serialized under, and can be retrieved by
providing, the specified key string.
*/
type DataAccess interface {
	Save(table, key string, object interface{}) error
	Read(table, key string, opts QueryOptions, object interface{}) error
	ReadAll(table string, readType ReadAllInterface) ([]interface{}, error)
	FilteredReadAll(table string, opts QueryOptions,
		readType ReadAllInterface) ([]interface{}, error)
	Delete(table, id string, opts QueryOptions, objects ...interface{}) error
}

/*
ReadAllInterface implementations must contain a GetType() function that returns
the type of the implementor.
*/
type ReadAllInterface interface {
	GetType() interface{}
}

/*
Query options are passed to the DataAccess layer from a controller
*/
type QueryOptions struct {
	Filters []Filter
}
/*
NewQueryOptions creates a new options object.
key: query field name
value: query value
id: True if field is a UUID key
*/
func NewQueryOptions(key, value string, id bool) QueryOptions {
	filter := Filter{key, value, id}
	return QueryOptions{
		Filters: []Filter{filter},
	}
}
/*
Filters are additional arguments used to narrow the number of objects returned
or affected by a particular Query
*/
type Filter struct {
	key		string
	value	interface{}
	id		bool
}
/*
AddFilter adds a filter to an QueryOptions object
key: query field name
value: query value
id: True if field is a UUID key
*/
func (o *QueryOptions) AddFilter(key, value string, id bool) {
	filter := Filter{key, value, id}
	o.Filters = append(o.Filters, filter)
}
