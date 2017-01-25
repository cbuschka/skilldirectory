package data

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// LocalFileSystem represents the project's directory on the local machine's
// file system. Implements the data.FileSystem interface.
type LocalFileSystem struct {
	protocol string
	hostname string
	rootdir  string
}

// NewLocalFileSystem returns a new LocalFileSystem object initialized to
// operate within the project's directory in the local file system
func NewLocalFileSystem() *LocalFileSystem {
	return &LocalFileSystem{
		protocol: "file://",
		hostname: "/",
		rootdir:  "skilldirectory/",
	}
}

// Read returns an io.Reader tht reads from the resource at the specified path
// within the project's directory in the local file system.
func (lfs *LocalFileSystem) Read(path string) (resource io.Reader, err error) {
	fullPath := lfs.hostname + lfs.rootdir + path
	// Open file on local file system (return error if fails or file doesn't exist)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %q :%s", fullPath, err)
	}

	// Extract data from file and close
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %q: %s", fullPath, err)
	}
	file.Close()

	// Return file's data
	return bytes.NewReader(fileBytes), nil // Read successful!
}

// Write saves the specified resource to the project's directory on the local
// file system under the specified path.
func (lfs *LocalFileSystem) Write(path string, resource io.ReadSeeker) (url string,
	err error) {
	fullPath := lfs.hostname + lfs.rootdir + path
	// Create file on local file system (or truncate and open if already exists)
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %q: %s", fullPath, err)
	}

	// Extract data from passed-in resource and write to local file system
	resourceBytes, _ := ioutil.ReadAll(resource)
	if len(resourceBytes) == 0 {
		return "", fmt.Errorf("please pass in a resource with > 0 bytes to write")
	}
	_, err = file.Write(resourceBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write data to file: %q: %s", fullPath, err)
	}

	// Successfully wrote resource to disk!
	url = lfs.protocol + fullPath
	return url, nil
}

// Delete removes the resource located at the specified path from the project's
// directory in the local file system.
func (lfs *LocalFileSystem) Delete(path string) (err error) {
	fullPath := lfs.hostname + lfs.rootdir + path
	// Delete file from local file system
	err = os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %q: %s", fullPath, err)
	}
	return nil // Delete successful!
}
