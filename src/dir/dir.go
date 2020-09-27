package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const InitialCapacity = 1 << 7

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
	dir := New(reader.DirPath(), InitialCapacity)

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
	dir := New(reader.DirPath(), InitialCapacity)

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

// CleanEmptyDir recursively cleans up empty sub catalogs in the tree.
// It might be especially useful to run this method on the diff.
func (tree *Dir) CleanEmptyDir() {
	for dirName, subTree := range tree.SubDirs {
		if subTree.IsEmpty() {
			delete(tree.SubDirs, dirName)
			continue
		}
		subTree.CleanEmptyDir()
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
	const indent = 4

	if depth > 0 {
		depth++
	}

	addSpaces(s, depth*indent)
	s.WriteString(fmt.Sprintf("[%s]:\n", tree.Path))
	for _, file := range tree.Files {
		addSpaces(s, (depth+1)*indent)
		s.WriteString(fmt.Sprintf("-%s (%d)\n", file.Name(), file.Size()))
	}

	for dirName, subDir := range tree.SubDirs {
		addSpaces(s, (depth+1)*indent)
		s.WriteString(fmt.Sprintf("-%s (DIR)\n", dirName))

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
