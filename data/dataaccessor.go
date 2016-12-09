package data

type DataAccessor struct {
	DataAccess
}

func NewAccessor(access DataAccess) DataAccessor {
	return DataAccessor{access}
}

type DataAccess interface {
	Save(string, interface{}) error
	Read(string, interface{}) error
	ReadAll(string, ReadAllInterface) ([]interface{}, error)
	FilteredReadAll(string, ReadAllInterface, func (interface{}) bool) ([]interface{}, error)
	Delete(string) error
}

type ReadAllInterface interface {
	GetType()interface{}
	// GetSlice()[]interface{}
}
