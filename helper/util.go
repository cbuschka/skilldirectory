package helper

import "strings"

// Returns the directory at the root of the specified path. Ignores starting slashes (regards
// "/skills/files" as "skills/files". Calling with "skills/files/whatever/1234-5678-9101" would return "skills/".
func getRootDir(path string) string {
	if path[0] == '/' {
		path = path[1:]
	}
	var rootDir string
	if strings.Index(path, "/") != -1 {
		rootDir = path[:strings.Index(path, "/")+1]
	} else {
		rootDir = path + "/"
	}
	return rootDir
}
