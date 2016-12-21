package data

import (
	"fmt"
)

type MockDataAccessor struct{}

func (m MockDataAccessor) Save(t, s string, i interface{}) error { return nil }
func (m MockDataAccessor) Read(t, s string, i interface{}) error { return nil }
func (m MockDataAccessor) Delete(t, s string) error              { return nil }
func (m MockDataAccessor) ReadAll(t string, r ReadAllInterface) ([]interface{}, error) {
	return nil, nil
}
func (d MockDataAccessor) FilteredReadAll(t string, opts Options, r ReadAllInterface) ([]interface{}, error) {
	return nil, nil
}

type MockErrorDataAccessor struct{}

func (e MockErrorDataAccessor) Save(t, s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Read(t, s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Delete(t, s string) error              { return fmt.Errorf("") }
func (e MockErrorDataAccessor) ReadAll(t string, r ReadAllInterface) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}
func (d MockErrorDataAccessor) FilteredReadAll(t string, opts Options, r ReadAllInterface) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////// MockInMemoryDataAccessor /////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
MockInMemoryDataAccessor implements the dataccessor.DataAccess interface by providing in-memory data storage,
designed to facilitate easy testing of components that rely on skill storage and retrieval without requiring access to
an external filesystem or database when running tests. Because MockInMemoryDataAccessor uses a computer's volatile
memory, it should not be used for permanent data storage, and is unlikely to have a use outside of unit testing.
*/
// type MockInMemoryDataAccessor struct {
// 	dataMap map[string][]byte // map of object to byte slice
// }
//
// func NewMockInMemoryDataAccessor() MockInMemoryDataAccessor {
// 	var retVal MockInMemoryDataAccessor
// 	retVal.dataMap = make(map[string][]byte)
// 	return retVal
// }
//
// func (e MockInMemoryDataAccessor) Save(table, ID string, object interface{}) error {
// 	b, err := json.Marshal(object)
// 	if err != nil {
// 		return err
// 	}
// 	e.dataMap[ID] = b
// 	return nil
// }
//
// func (e MockInMemoryDataAccessor) Read(table, ID string, object interface{}) error {
// 	data := e.dataMap[ID]
// 	if len(data) == 0 {
// 		return fmt.Errorf("No such object with ID: %s", ID)
// 	}
// 	json.Unmarshal(data, &object)
// 	return nil
// }
//
// func (e MockInMemoryDataAccessor) Delete(table, ID string) error {
// 	fmt.Println("Deleting this key from map:", ID)
// 	data := e.dataMap[ID]
// 	if len(data) == 0 {
// 		return fmt.Errorf("No such object with ID: %s", ID)
// 	}
// 	e.dataMap[ID] = make([]byte, 0)
// 	return nil
// }
//
// func (e MockInMemoryDataAccessor) ReadAll(table, path string, readType ReadAllInterface) ([]interface{}, error) {
// 	returnObjects := []interface{}{}
// 	object := readType.GetType()
// 	for _, val := range e.dataMap {
// 		json.Unmarshal(val, object)
// 		returnObjects = append(returnObjects, object)
// 	}
// 	return returnObjects, nil
// }
//
// func (e MockInMemoryDataAccessor) FilteredReadAll(table, path string, opts Options, readType ReadAllInterface) ([]interface{}, error) {
// 	returnObjects := []interface{}{}
// 	object := readType.GetType()
// 	for _, val := range e.dataMap {
// 		json.Unmarshal(val, &object)
// 			returnObjects = append(returnObjects, object)
// 		}
// 	}
// 	return returnObjects, nil
// }

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
