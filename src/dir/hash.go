package dir

import (
	"bytes"
	"crypto/md5"
)

// Hash method calculates MD5 hash based on names of all files and path in
// given tree.
func (tree *Dir) Hash() [md5.Size]byte {
	var b bytes.Buffer
	treeNames := tree.forHash(&b)
	return md5.Sum(treeNames.Bytes())
}

func (tree *Dir) forHash(b *bytes.Buffer) *bytes.Buffer {
	b.WriteString(tree.Path)

	for _, file := range tree.Files {
		b.WriteString(file.Name())
	}

	for _, subDir := range tree.SubDirs {
		if subDir.Path != "" {
			subDir.forHash(b)
		}
	}

	return b
}
