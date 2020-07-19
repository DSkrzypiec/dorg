package dir

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Dir represents tree structure of files and catalogs.
type Dir struct {
	Path    string
	Files   []os.FileInfo
	SubDirs map[string]Dir
}

// DirScanConfig contains configuration for scanning file tree.
type DirScanConfig struct {
	ScanOnlyTopLevel bool
	ExcludedDirs     map[string]struct{}
	ScanTimeout      time.Duration
}

type FileTreeScanner interface {
	Scan(string, DirScanConfig) (Dir, error)
}

// Scanner represents a scanner of OS file system.
type Scanner struct{}

// CurrentDirScanner represents a scanner of OS file system but without going
// into sub catalogs.
type CurrentDirScanner struct{}

// TODO
func (s Scanner) Scan(path string, config DirScanConfig) (Dir, error) {
	return Dir{}, nil
}

// Scan scans files and sub catalogs from path
func (s CurrentDirScanner) Scan(path string, config DirScanConfig) (Dir, error) {
	currentDir, err := os.Open(path)
	if err != nil {
		return Dir{}, errors.Wrap(err, fmt.Sprintf("Cannot open path [%s]", path))
	}

	fileInfos, err := currentDir.Readdir(-1)
	if err != nil {
		return Dir{}, errors.Wrap(err, fmt.Sprintf("Cannot read file info for all files from [%s]", path))
	}

	files := make([]os.FileInfo, 0, 100)
	subDirs := make(map[string]Dir)

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		_, isExcluded := config.ExcludedDirs[name]

		if fileInfo.IsDir() && !isExcluded {
			subDirs[name] = Dir{}
			continue
		}

		if !fileInfo.IsDir() {
			files = append(files, fileInfo)
		}
	}

	return Dir{path, files, subDirs}, nil
}

func (tree Dir) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("Files in [%s]:\n", tree.Path))
	for _, file := range tree.Files {
		s.WriteString(fmt.Sprintf("  -%s (%d)\n", file.Name(), file.Size()))
	}

	s.WriteString(fmt.Sprintf("Catalogs in [%s]:\n", tree.Path))
	for dirName, _ := range tree.SubDirs {
		s.WriteString(fmt.Sprintf("  -%s\n", dirName))
	}

	return s.String()
}
