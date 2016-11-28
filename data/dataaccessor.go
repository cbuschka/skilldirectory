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
	Delete(string) error
}
