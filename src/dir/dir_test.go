package dir

import (
	"testing"
)

func TestScanTopLevel(t *testing.T) {
	config := DirScanConfig{}

	tree, err := ScanTopLevel("../", config)
	if err != nil {
		t.Errorf("Error while scanning dir tree [../]: %s", err.Error())
	}

	if len(tree.Files) <= 2 {
		t.Errorf("Expected at least 3 files in [src/], got [%d]", len(tree.Files))
	}
	if len(tree.SubDirs) <= 1 {
		t.Errorf("Expected at least 2 sub catalogs of [src/], got: [%d]", len(tree.SubDirs))
	}
	if _, ok := tree.SubDirs["config"]; !ok {
		t.Error("Missing expected [../config] catalog in the tree.")
	}
	if _, ok := tree.SubDirs["dir"]; !ok {
		t.Error("Missing expected [../dir] catalog in the tree.")
	}
}

func TestScan(t *testing.T) {
	config := DirScanConfig{}

	tree, err := Scan("../", config)
	if err != nil {
		t.Errorf("Error while scanning dir tree [../]: %s", err.Error())
	}

	if len(tree.Files) <= 2 {
		t.Errorf("Expected at least 3 files in [src/], got [%d]", len(tree.Files))
	}
	if len(tree.SubDirs) <= 1 {
		t.Errorf("Expected at least 2 sub catalogs of [src/], got: [%d]", len(tree.SubDirs))
	}
	if _, ok := tree.SubDirs["config"]; !ok {
		t.Error("Missing expected [../config] catalog in the tree.")
	}
	if _, ok := tree.SubDirs["dir"]; !ok {
		t.Error("Missing expected [../dir] catalog in the tree.")
	}

	dirTree := tree.SubDirs["dir"]
	if len(dirTree.Files) <= 1 {
		t.Errorf("Expected at least 2 files in catalog [src/dir/], got: [%d]", len(dirTree.Files))
	}
}

func TestHash(t *testing.T) {
	config := DirScanConfig{}

	tree, err := Scan("../", config)
	if err != nil {
		t.Errorf("Error while scanning dir tree [../]: %s", err.Error())
	}

	dirTree, dirErr := Scan("../dir", config)
	if dirErr != nil {
		t.Errorf("Error while scanning dir tree [../dir/]: %s", dirErr.Error())
	}

	subDir, exists := tree.SubDirs["dir"]
	if !exists {
		t.Errorf("Missing expected [../dir] catalog in the tree.")
	}

	h1 := subDir.Hash()
	h2 := dirTree.Hash()

	hashesAreTheSame := true
	for id, _ := range h1 {
		if h1[id] != h2[id] {
			hashesAreTheSame = false
			break
		}
	}

	if !hashesAreTheSame {
		t.Errorf("Expected the same hash for [../dir], got \n [%v] vs \n [%v]",
			h1, h2)
	}
}

func TestEquals(t *testing.T) {
	config := DirScanConfig{}

	tree, err := Scan("../", config)
	if err != nil {
		t.Errorf("Error while scanning dir tree [../]: %s", err.Error())
	}

	dirTree, dirErr := Scan("../dir", config)
	if dirErr != nil {
		t.Errorf("Error while scanning dir tree [../dir/]: %s", dirErr.Error())
	}

	subDir, exists := tree.SubDirs["dir"]
	if !exists {
		t.Errorf("Missing expected [../dir] catalog in the tree.")
	}

	if !dirTree.Equals(subDir) {
		t.Errorf("Subtree [dir] of [../] suppose to be the same as [../dir].")
	}
}
