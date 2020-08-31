package dir

import (
	"os"
	"time"
)

type MockDir struct {
	Path    string
	Files   []os.FileInfo
	SubDirs map[string]MockDir
}

func (md MockDir) DirPath() string {
	return md.Path
}

func (md MockDir) Readdir() ([]os.FileInfo, error) {
	return md.Files, nil
}

func (md MockDir) New(path string) DirReader {
	if _, exist := md.SubDirs[path]; exist {
		return md.SubDirs[path]
	}

	emptyDirs := make(map[string]MockDir)
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
