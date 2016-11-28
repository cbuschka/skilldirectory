package data

import "fmt"

type DataMock struct{}

func (d DataMock) Read(key string, object interface{}) error {
	fmt.Println("Reading")
	return nil
}
func (d DataMock) Save(key string, object interface{}) error {
	fmt.Println("Writing")
	return nil
}
func (d DataMock) Delete(key string) error {
	fmt.Println("Deleting")
	return nil
}
