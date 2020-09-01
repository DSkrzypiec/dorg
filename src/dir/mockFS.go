package dir

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MockDir struct {
	Path    string
	Files   []os.FileInfo
	SubDirs map[string]*MockDir
}

func (md MockDir) DirPath() string {
	return md.Path
}

func (md MockDir) Readdir() ([]os.FileInfo, error) {
	files := make([]os.FileInfo, 0, 10)

	for _, file := range md.Files {
		files = append(files, file)
	}

	for dirName := range md.SubDirs {
		files = append(files, NewMockFileInfo(dirName, true))
	}

	return files, nil
}

func (md MockDir) New(path string) DirReader {
	for _, dir := range md.SubDirs {
		if dir.Path == path {
			return dir
		}
	}

	emptyDirs := make(map[string]*MockDir)
	return MockDir{path, make([]os.FileInfo, 0), emptyDirs}
}

type MockFileInfo struct {
	FileName    string
	IsDirectory bool
}

func (mfi MockFileInfo) Name() string       { return mfi.FileName }
func (mfi MockFileInfo) Size() int64        { return int64(8) }
func (mfi MockFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (mfi MockFileInfo) ModTime() time.Time { return time.Now() }
func (mfi MockFileInfo) IsDir() bool        { return mfi.IsDirectory }
func (mfi MockFileInfo) Sys() interface{}   { return nil }

func NewMockFileInfo(name string, isDir bool) MockFileInfo {
	return MockFileInfo{
		FileName:    name,
		IsDirectory: isDir,
	}
}

// Creates new MockDir. For convenience only files are required but the
// following convention is assumed:
//	* if element of files contains "_" e.g. "xyz_file1.go" then file1.go will
//	be put inside "xyz" sub directory.
// * if element of files doesn't contain "_" then it will be put in top level
// directory.
// In order to produce deeper trees use this function several times for
// sub-trees.
func NewMockDir(rootPath string, files ...string) MockDir {
	topFiles := make([]os.FileInfo, 0, 10)
	subDirs := make(map[string]*MockDir)

	for _, file := range files {
		if strings.Contains(file, "_") {
			parts := strings.Split(file, "_")
			dirName := parts[0]
			subFile := NewMockFileInfo(parts[1], false)

			if _, exist := subDirs[dirName]; exist {
				(*subDirs[dirName]).Files = append((*subDirs[dirName]).Files, subFile)
				continue
			}

			subDirs[dirName] = &MockDir{filepath.Join(rootPath, dirName), []os.FileInfo{subFile}, nil}
			continue
		}
		topFiles = append(topFiles, NewMockFileInfo(file, false))
	}

	return MockDir{rootPath, topFiles, subDirs}
}
