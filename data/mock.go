package data

import (
	"fmt"
	"io"
)

type MockDataAccessor struct{}

func (m MockDataAccessor) Save(t, s string, i interface{}) error { return nil }
func (m MockDataAccessor) Read(t, s string, opts QueryOptions, i interface{}) error {
	return nil
}
func (c MockDataAccessor) Delete(table, id string, opts QueryOptions, objects ...interface{}) error { return nil }
func (m MockDataAccessor) ReadAll(t string, r ReadAllInterface) ([]interface{}, error) {
	return nil, nil
}
func (d MockDataAccessor) FilteredReadAll(t string, opts QueryOptions, r ReadAllInterface) ([]interface{}, error) {
	return nil, nil
}

type MockErrorDataAccessor struct{}

func (e MockErrorDataAccessor) Save(t, s string, i interface{}) error { return fmt.Errorf("") }
func (e MockErrorDataAccessor) Read(t, s string, opts QueryOptions, i interface{}) error {
	return fmt.Errorf("")
}
func (c MockErrorDataAccessor) Delete(table, id string, opts QueryOptions, objects ...interface{}) error {
	return fmt.Errorf("")
}
func (e MockErrorDataAccessor) ReadAll(t string, r ReadAllInterface) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}
func (d MockErrorDataAccessor) FilteredReadAll(t string, opts QueryOptions, r ReadAllInterface) ([]interface{}, error) {
	return nil, fmt.Errorf("")
}

type MockFileSystem struct{}

func (m MockFileSystem) Read(path string) (resource io.Reader, err error) {
	return nil, nil
}

func (m MockFileSystem) Write(path string, resource io.ReadSeeker) (url string,
	err error) {
	return "", nil
}

func (m MockFileSystem) Delete(path string) (err error) {
	return nil
}

type MockErrorFileSystem struct{}

func (m MockErrorFileSystem) Read(path string) (resource io.Reader, err error) {
	return nil, fmt.Errorf("")
}

func (m MockErrorFileSystem) Write(path string, resource io.ReadSeeker) (url string,
	err error) {
	return "", fmt.Errorf("")
}

func (m MockErrorFileSystem) Delete(path string) (err error) {
	return fmt.Errorf("")
}
