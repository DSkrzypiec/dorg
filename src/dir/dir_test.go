package dir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanTopLevel(t *testing.T) {
	config := DirScanConfig{}

	tree, err := ScanTopLevel(FlatDir{"../"}, config)
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

	tree, err := Scan(FlatDir{"../"}, config)
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

func TestScanExclude(t *testing.T) {
	excludedDirs := make(map[string]struct{})
	var empty struct{}
	excludedDirs["sub1"] = empty
	config := DirScanConfig{ExcludedDirs: excludedDirs}

	md := NewMockDir("Downloads", "f1.go", "f2.cpp", "sub1_g1.exe", "sub1_g2.cs", "xyz_h1.jpeg")
	tree, err := Scan(md, config)
	if err != nil {
		t.Errorf("Error while scanning mock dir tree [Downloads]: %s", err.Error())
	}

	if len(tree.Files) != 2 {
		t.Errorf("Excpected 2 files in top dir, got: %d", len(tree.Files))
	}
	if len(tree.SubDirs) != 1 {
		t.Errorf("Excepted only one subcatalog, got: %d", len(tree.SubDirs))
	}
	if _, exist := tree.SubDirs["sub1"]; exist {
		t.Error("Catalog 'sub1' suppose to be excluded but wasn't")
	}
}

func TestScanDeep(t *testing.T) {
	config := DirScanConfig{}
	md := NewMockDir("Downloads", "f1.go", "f2.cpp", "sub1_g1.exe", "sub1_g2.cs", "xyz_h1.jpeg")
	subMd := NewMockDir(filepath.Join("Downloads", "sub2"), "d2.a", "d22.b", "sub2sub_e0", "sub2sub_e10")
	md.SubDirs["sub2"] = &subMd

	tree, err := Scan(md, config)
	if err != nil {
		t.Errorf("Error while scanning mock dir tree [Downloads]: %s", err.Error())
	}

	//fmt.Println(tree)

	if len(tree.Files) != 2 {
		t.Errorf("Expected exactly 2 files in root dir. Got %d files - [%v]", len(tree.Files), tree.Files)
	}
	if len(tree.SubDirs) != 3 {
		t.Errorf("Expected exactly 3 sub-directories. Got %d dirs", len(tree.SubDirs))
	}
	if _, exist := tree.SubDirs["sub1"]; !exist {
		t.Errorf("Missing expected 'sub1' sub directory")
	}
	if _, exist := tree.SubDirs["xyz"]; !exist {
		t.Errorf("Missing expected 'xyz' sub directory")
	}
	if _, exist := tree.SubDirs["sub2"]; !exist {
		t.Errorf("Missing expected 'sub2' sub directory")
	}

	sub2 := tree.SubDirs["sub2"]
	if len(sub2.Files) != 2 {
		t.Errorf("Expected exactly 2 files in 'sub2' sub directory. Got %d files - [%v]",
			len(sub2.Files), sub2.Files)
	}
	if len(sub2.SubDirs) != 1 {
		t.Errorf("Expected single sub directory in 'Downloads/sub2'. Got %d",
			len(sub2.SubDirs))
	}
}

func TestCleanEmptyDirs(t *testing.T) {
	// Arrange
	tree3 := Dir{
		Path:  "Path1/sub3",
		Files: []os.FileInfo{NewMockFileInfo("g.go", false)},
		SubDirs: map[string]Dir{
			"x1": New("Path1/sub3/x1", 1),
			"x2": Dir{
				Path:  "Path1/sub3/x2",
				Files: []os.FileInfo{NewMockFileInfo("main.cpp", false)},
			},
		},
	}

	subDirs1 := make(map[string]Dir)
	subDirs1["sub1"] = New("Path1/sub1", 1)
	subDirs1["sub2"] = New("Path1/sub2", 1)
	subDirs1["sub3"] = tree3

	tree := Dir{
		Path:    "Path1",
		Files:   []os.FileInfo{NewMockFileInfo("f1.txt", false)},
		SubDirs: subDirs1,
	}

	// Act
	tree.CleanEmptyDir()

	// Assert
	if len(tree.Files) != 1 {
		t.Errorf("Number of files in root dir should be 1, got: %d",
			len(tree.Files))
	}
	if len(tree.SubDirs) != 1 {
		t.Errorf("Expected 1 dir after clean up, got: %d", len(tree.SubDirs))
	}

	sub3, exist := tree.SubDirs["sub3"]
	if !exist {
		t.Error("Sub catalog [sub3] should not be deleted.")
	}
	if len(sub3.Files) != 1 {
		t.Errorf("Number of files in [Path1/sub3] should be 1, got: %d",
			len(sub3.Files))
	}
	if len(sub3.SubDirs) != 1 {
		t.Errorf("Expected 1 dir after clean up in [Path1/sub3], got: %d",
			len(sub3.SubDirs))
	}
	if _, x2Exist := sub3.SubDirs["x2"]; !x2Exist {
		t.Error("Sub catalog [x2] should exist in [Path1/sub3] after cleanup.")
	}
}

func TestHash(t *testing.T) {
	config := DirScanConfig{}

	tree, err := Scan(FlatDir{".."}, config)
	if err != nil {
		t.Errorf("Error while scanning dir tree [../]: %s", err.Error())
	}

	dirTree, dirErr := Scan(FlatDir{filepath.Join("..", "dir")}, config)
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

	tree, err := Scan(FlatDir{".."}, config)
	if err != nil {
		t.Errorf("Error while scanning dir tree [../]: %s", err.Error())
	}

	dirTree, dirErr := Scan(FlatDir{filepath.Join("..", "dir")}, config)
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

func TestDirFilesCount(t *testing.T) {
	config := DirScanConfig{}
	md := NewMockDir("Downloads", "f1.go", "f2.cpp", "sub1_g1.exe", "sub1_g2.cs",
		"xyz_h1.jpeg")
	subMd := NewMockDir(filepath.Join("Downloads", "sub2"), "d2.a", "d22.b",
		"sub2sub_e0", "sub2sub_e10")
	md.SubDirs["sub2"] = &subMd

	tree, err := Scan(md, config)
	if err != nil {
		t.Errorf("Error while scanning mock dir tree [Downloads]: %s", err.Error())
	}

	fCnt := tree.FilesCount()
	if fCnt != 9 {
		t.Errorf("Tree has 9 files, but counted: %d", fCnt)
	}

	sub2Tree, exists := tree.SubDirs["sub2"]
	if !exists {
		t.Errorf("Subcatalog 'sub2' doesn't exist in the tree but suppose to")
	}

	sub2FileCnt := sub2Tree.FilesCount()
	if sub2FileCnt != 4 {
		t.Errorf("Subcatalog 'sub2' has 4 files, but counted: %d", sub2FileCnt)
	}

	empty := New("path", 100)
	emptyCnt := empty.FilesCount()
	if emptyCnt != 0 {
		t.Errorf("Empty Dir suppose to have 0 files, but counted: %d", emptyCnt)
	}
}
