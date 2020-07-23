package dir

// Checks if given tree is identical as the object.
func (tree Dir) Equals(another Dir) bool {
	hash := tree.Hash()
	anotherHash := another.Hash()

	for id, _ := range hash {
		if hash[id] != anotherHash[id] {
			return false
		}
	}
	return true
}

// Diff compares given (another) file tree with current file tree. First return
// states if there is a difference between trees and another produce "diff
// tree". This method is meant to be used to check if particular file tree
// changes in time. For Dirs with different root Paths result is
// (true, tree \ another). Moreover in that case difference tree \ another may
// don't have any sense.
func (tree Dir) Diff(another Dir) (bool, Dir) {
	// TODO
	return true, tree
}
