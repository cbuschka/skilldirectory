package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
