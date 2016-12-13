package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Implments DataAccessor Type
type FileWriter struct {
	path string
}

func NewFileWriter(path string) FileWriter {
	return FileWriter{path}
}

func (f FileWriter) Save(key string, object interface{}) error {
	path := f.pathHelper(key)
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0600)
}

func (f FileWriter) Read(key string, object interface{}) error {
	path := f.pathHelper(key)
	data, err := ioutil.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &object)
	}
	return err
}

func (f FileWriter) Delete(key string) error {
	path := f.pathHelper(key)
	return os.Remove(path)

}

func (f FileWriter) pathHelper(key string) string {
	return f.path + key
}

func (f FileWriter) ReadAll(path string, readType ReadAllInterface) ([]interface{}, error) {
	returnObjects := []interface{}{}
	object := readType.GetType()
	filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			err1 := f.Read(fi.Name(), &object)
			if err1 != nil {
				return err1
			}
			returnObjects = append(returnObjects, object)
		}
		return nil
	})
	return returnObjects, nil
}

/*
Applies the specified filterFunc to each read from database/repository. Returns a slice containing entries that satisfied
the filtering function.
*/
func (f FileWriter) FilteredReadAll(path string, readType ReadAllInterface,
	filterFunc func(interface{}) bool) ([]interface{}, error) {
	returnObjects := []interface{}{}
	object := readType.GetType()
	filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			err1 := f.Read(fi.Name(), &object)
			if err1 != nil {
				return err1
			}
			if filterFunc(object) {
				returnObjects = append(returnObjects, object)
			}
		}
		return nil
	})
	return returnObjects, nil
}
