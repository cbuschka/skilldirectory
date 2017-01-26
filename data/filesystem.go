package data

import (
	"io"
)

/*
FileSystem represents a file system that contains resources identifable by a
path string. These resources can be Read, Written, or Deleted.
*/
type FileSystem interface {
	Read(path string) (resource io.Reader, err error)
	Write(path string, resource io.ReadSeeker) (url string, err error)
	Delete(path string) (err error)
}
