package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Dir represents tree structure of files and catalogs.
type Dir struct {
	Path    string
	Files   []os.FileInfo
	SubDirs map[string]Dir
}

// DirScanConfig contains configuration for scanning file tree.
type DirScanConfig struct {
	ExcludedDirs map[string]struct{}
}

// Scan scans whole OS file system tree starting from path.
func Scan(reader DirReader, config DirScanConfig) (Dir, error) {
	fileInfos, err := reader.Readdir()
	if err != nil {
		return Dir{}, err
	}
	dir := New(reader.DirPath(), 100)

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		if _, isExcluded := config.ExcludedDirs[name]; isExcluded {
			continue
		}

		if !fileInfo.IsDir() {
			dir.Files = append(dir.Files, fileInfo)
			continue
		}

		newPath := filepath.Join(reader.DirPath(), name)
		newReader := reader.New(newPath)
		subDir, err := Scan(newReader, config)
		if err != nil {
			return Dir{}, err
		}
		dir.SubDirs[name] = subDir
	}

	return dir, nil
}

// ScanTopLevel scans files and sub catalogs from path without recurring sub catalogs.
func ScanTopLevel(reader DirReader, config DirScanConfig) (Dir, error) {
	fileInfos, err := reader.Readdir()
	if err != nil {
		return Dir{}, err
	}
	dir := New(reader.DirPath(), 100)

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		_, isExcluded := config.ExcludedDirs[name]

		if fileInfo.IsDir() && !isExcluded {
			dir.SubDirs[name] = Dir{}
			continue
		}

		if !fileInfo.IsDir() {
			dir.Files = append(dir.Files, fileInfo)
		}
	}

	return dir, nil
}

// Initialize new Dir.
func New(path string, capacity int) Dir {
	return Dir{
		Path:    path,
		Files:   make([]os.FileInfo, 0, capacity),
		SubDirs: make(map[string]Dir),
	}
}

// Checks whenever Dir is empty.
func (tree Dir) IsEmpty() bool {
	return tree.FilesCount() == 0
}

// Return number of files in Dir (including all sub-catalogs recursively).
func (tree Dir) FilesCount() int {
	nFiles := 0
	filesCount(tree, &nFiles)
	return nFiles
}

// Recursive files count.
func filesCount(tree Dir, currentCount *int) {
	*currentCount += len(tree.Files)

	for _, subTree := range tree.SubDirs {
		filesCount(subTree, currentCount)
	}
}

// Tree in form of formatted string.
func (tree Dir) String() string {
	var s strings.Builder
	return tree.print(&s, 0)
}

// Recursion tree printing helper.
func (tree Dir) print(s *strings.Builder, depth int) string {
	const indent = 2

	addSpaces(s, depth*indent)
	s.WriteString(fmt.Sprintf("Files in [%s]:\n", tree.Path))
	for _, file := range tree.Files {
		addSpaces(s, depth*indent)
		s.WriteString(fmt.Sprintf("  -%s (%d)\n", file.Name(), file.Size()))
	}

	addSpaces(s, depth*indent)
	s.WriteString(fmt.Sprintf("Catalogs in [%s]:\n", tree.Path))
	for dirName, subDir := range tree.SubDirs {
		addSpaces(s, depth*indent)
		s.WriteString(fmt.Sprintf("  -%s\n", dirName))

		if subDir.Path != "" {
			subDir.print(s, depth+1)
		}
	}

	return s.String()
}

func addSpaces(builder *strings.Builder, n int) *strings.Builder {
	for i := 0; i < n; i++ {
		builder.WriteString(" ")
	}
	return builder
}
