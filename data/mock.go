package data

import (
	"fmt"
	"io"
)

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
