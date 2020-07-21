package dir

// Checks if given tree is identical as the object.
func (tree *Dir) Equals(another *Dir) bool {
	hash := tree.Hash()
	anotherHash := another.Hash()

	for id, _ := range hash {
		if hash[id] != anotherHash[id] {
			return false
		}
	}
	return true
}

// TODO: Other function for selecting difference between two trees
