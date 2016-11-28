package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DataAccess interface {
	Save(string, string, interface{}) error
	Read(string, string, interface{}) error
}

type FileWriter struct{}

func (f FileWriter) Save(pathName, fileName string, object interface{}) error {
	path := pathHelper(pathName, fileName)
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0600)
}

func (f FileWriter) Read(pathName, fileName string, object interface{}) error {
	path := pathHelper(pathName, fileName)
	data, err := ioutil.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &object)
	}
	return err
}

func (f FileWriter) Delete(pathName, fileName string) error {
	path := pathHelper(pathName, fileName)
	return os.Remove(path)

}

func pathHelper(pathName, name string) string {
	return pathName + "/" + name + ".txt"
}
