package dir

import (
	"os"
	"path/filepath"
)

// Checks if given tree is identical as the object.
func (tree Dir) Equals(another Dir) bool {
	hash := tree.Hash()
	anotherHash := another.Hash()

	for id := range hash {
		if hash[id] != anotherHash[id] {
			return false
		}
	}
	return true
}

// Diff compares given (another) file tree with current file tree. First return
// states if there is a difference between trees and another produce "diff
// tree" which is tree \ another in term of set subtraction. This method is
// meant to check whenever given Dir have changed over some time period. In
// particular if tree.Path and another.Path differs then Diff return false and
// empty Dir. It wouldn't return true and actual difference between trees.
func (tree Dir) Diff(another Dir) (bool, Dir) {
	if tree.Path != another.Path {
		return false, Dir{}
	}

	diffTree := New(tree.Path, InitialCapacity)
	dirDiff(tree, another, &diffTree)

	return !diffTree.IsEmpty(), diffTree
}

func dirDiff(tree, another Dir, diff *Dir) Dir {
	diff.Files = filesDiff(tree.Files, another.Files)

	for dirName, sub := range tree.SubDirs {
		anotherSub, existsInAnother := another.SubDirs[dirName]
		if !existsInAnother {
			diff.SubDirs[dirName] = sub
			continue
		}

		newSub := New(filepath.Join(tree.Path, dirName), InitialCapacity)
		diff.SubDirs[dirName] = dirDiff(sub, anotherSub, &newSub)
	}

	return *diff
}

// Determines orig \ new.
func filesDiff(orig, new []os.FileInfo) []os.FileInfo {
	diff := make([]os.FileInfo, 0, len(orig))
	newMap := make(map[string]os.FileInfo)

	for _, file := range new {
		if _, exists := newMap[file.Name()]; !exists {
			newMap[file.Name()] = file
		}
	}

	for _, origFile := range orig {
		if _, existInNew := newMap[origFile.Name()]; !existInNew {
			diff = append(diff, origFile)
		}
	}

	return diff
}
