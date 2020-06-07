// Package ddir contains representation of downloads directory
// (ddir) and functionality around it.
package ddir

import (
	"os"
)

// DDir represents tree structure of files and catalogs inside so-called
// downloads directory.
type DDir struct {
	Files []*os.FileInfo
	Dirs  map[string]DDir
}

func Scan(ddirPath string) (DDir, error) {
	// TODO
}

func Empty(capacity int) DDir {
	return DDir{
		Files: make([]*os.FileInfo, 0, capacity),
		Dirs:  make(map[string]DDir),
	}
}
