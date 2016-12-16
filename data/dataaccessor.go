package data

/*
DataAccessor is a wrapper type for implementations of the DataAccess interface.
*/
type DataAccessor struct {
	DataAccess
}

/*
NewAccessor() is a constructor for the DataAccessor type. It accepts implementations of the
DataAccess interface, and returns an instance of DataAccessor wrapping the passed-in DataAccess implementation.
*/
func NewAccessor(access DataAccess) DataAccessor {
	return DataAccessor{access}
}

/*
DataAccess represents an interface for accessing and managing a data storage system (e.g. a filesystem or database).
Implementations of DataAccess should ensure that the objects passed to them are properly serialized under, and can be
retrieved by providing, the specified key string.
*/
type DataAccess interface {
	Save(table, key string, object interface{}) error
	Read(table, key string, object interface{}) error
	ReadAll(table, path string, readType ReadAllInterface) ([]interface{}, error)
	FilteredReadAll(path string, readType ReadAllInterface, filterFunc func(interface{}) bool) ([]interface{}, error)
	Delete(table, key string) error
}

/*
ReadAllInterface implementations must contain a GetType() function that returns the type of the implementor.
*/
type ReadAllInterface interface {
	GetType() interface{}
}
