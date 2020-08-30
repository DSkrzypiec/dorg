package dir

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// DirReader represents an object which can reads os.FileInfo elements. Most
// commont case would be an file system directory reader where Readdir() tries
// to read its elements.
type DirReader interface {
	DirPath() string
	Readdir() ([]os.FileInfo, error)
	New(string) DirReader
}

// Implementation of DirReader interface. FlatDir represents single file system
// directory.
type FlatDir struct {
	Path string
}

// Returns directory path.
func (fd FlatDir) DirPath() string {
	return fd.Path
}

// Readdir reads all elements of Path directory as os.FileInfo slice. It's a
// wrapper on (*os.File).Readdir(int) method.
func (fd FlatDir) Readdir() ([]os.FileInfo, error) {
	return readdir(fd.Path)
}

// New creates new FlatDir for given path.
func (fd FlatDir) New(path string) DirReader {
	return FlatDir{Path: path}
}

// Reads element of given directory.
func readdir(path string) ([]os.FileInfo, error) {
	currentDir, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Cannot open path [%s]", path))
	}

	fileInfos, err := currentDir.Readdir(-1)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Cannot read file info for all files from [%s]", path))
	}

	return fileInfos, nil
}
