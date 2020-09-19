package dir

import (
	"os"
	"testing"
)

func TestDirDiff(t *testing.T) {
	config := DirScanConfig{}
	md1 := NewMockDir("Path1", "f1.go", "f2.cpp", "sub1_g1.go", "sub2_h1.html", "sub3_i3.jpeg")
	md2 := NewMockDir("Path1", "f1.go", "sub1_g1.go", "sub3_i3.jpeg", "sub3_i4.txt")

	tree1, err1 := Scan(md1, config)
	if err1 != nil {
		t.Errorf("Error while scanning first tree: %s", err1.Error())
	}

	tree2, err2 := Scan(md2, config)
	if err2 != nil {
		t.Errorf("Error while scanning second tree: %s", err2.Error())
	}

	// tree1 - tree1
	diff0, diffTree0 := tree1.Diff(tree1)
	if diff0 {
		t.Error("tree.Diff(tree) suppose to return false - there is no diff.")
	}
	if !diffTree0.IsEmpty() {
		t.Errorf("tree.Diff(tree) suppose to return empty files tree, got: \n %s",
			diffTree0.String())
	}

	// tree1 - tree2
	diff1, diffTree1 := tree1.Diff(tree2)
	if !diff1 {
		t.Error("tree1.Diff(tree2) suppose to return true, because tree1 is different then tree2.")
	}
	if diffTree1.FilesCount() != 2 {
		t.Errorf("tree1.Diff(tree2) should contain 2 files, got: %d \n %s",
			diffTree1.FilesCount(), diffTree1.String())
	}
	if diffTree1.Files[0].Name() != "f2.cpp" {
		t.Error("Tree tree1.Diff(tree2) doesn't have file 'f2.cpp' in top catalog")
	}

	diff2, diffTree2 := tree2.Diff(tree1)
	if !diff2 {
		t.Error("tree2.Diff(tree1) suppose to return true, because tree1 is different then tree2.")
	}
	if diffTree2.FilesCount() != 1 {
		t.Errorf("tree2.Diff(tree1) should contain 1 file, got: %d \n %s",
			diffTree2.FilesCount(), diffTree2.String())
	}
	if _, existSub3 := diffTree2.SubDirs["sub3"]; !existSub3 {
		t.Error("There should be subcatalog 'sub3' in diff tree tree2.Diff(tree1)")
	}

	sub3 := diffTree2.SubDirs["sub3"]
	if sub3.Files[0].Name() != "i4.txt" {
		t.Error("There should be file 'i4.txt' in subcatalog 'sub3' in diff tree tree2.Diff(tree1)")
	}
}

func TestDirDiffEmpty(t *testing.T) {
	config := DirScanConfig{}
	empty := New("Path1", 1)
	md := NewMockDir("Path1", "f1.go", "f2.cpp", "sub1_g1.go", "sub1_g2.html")

	tree, err := Scan(md, config)
	if err != nil {
		t.Errorf("Error while scanning tree: %s", err.Error())
	}

	diff0, diffTree0 := tree.Diff(empty)
	if !diff0 {
		t.Error("There should be difference between empty and non-empty tree")
	}
	if !diffTree0.Equals(tree) {
		t.Errorf("tree.Diff(empty) should be the same as tree. There isn't. Got tree \n %s \n and diff: \n %s",
			tree.String(), diffTree0.String())
	}

	_, diffTree1 := empty.Diff(tree)
	if !diffTree1.Equals(empty) {
		t.Errorf("empty.Diff(tree) should be the same as empty tree. There isn't. Got empty \n %s \n and diff: \n %s",
			empty.String(), diffTree1.String())
	}

	_, emptyDiff := empty.Diff(empty)
	if !empty.Equals(emptyDiff) {
		t.Errorf("empty.Diff(empty) should be empty, got: \n %s", emptyDiff.String())
	}
}

func TestFilesDiff(t *testing.T) {
	orig := []os.FileInfo{
		NewMockFileInfo("f1", false),
		NewMockFileInfo("f2", false),
		NewMockFileInfo("g3", false)}

	new := []os.FileInfo{
		NewMockFileInfo("f2", false),
		NewMockFileInfo("f1", false)}

	diff := filesDiff(orig, new)

	if len(diff) != 1 {
		t.Errorf("Expected diff to be length 1, got: %d", len(diff))
	}
	if diff[0].Name() != "g3" {
		t.Errorf("Expected 'g3', got: %s", diff[0].Name())
	}
	if diff[0].IsDir() {
		t.Errorf("Excepted to be regular file, not dir")
	}
}

func TestFilesDiffEmpty(t *testing.T) {
	orig := []os.FileInfo{
		NewMockFileInfo("f1", false),
		NewMockFileInfo("f2", false),
		NewMockFileInfo("g3", false)}

	emptyDiff := filesDiff(orig, orig)

	if len(emptyDiff) > 0 {
		t.Errorf("Expected to be empty, but got length = %d", len(emptyDiff))
	}
}

func TestFilesDiffFull(t *testing.T) {
	orig := []os.FileInfo{
		NewMockFileInfo("f1", false),
		NewMockFileInfo("f2", false),
		NewMockFileInfo("g3", false)}

	empty := []os.FileInfo{}
	fullDiff := filesDiff(orig, empty)

	if len(fullDiff) != len(orig) {
		t.Errorf("Expected length %d, got %d", len(orig), len(fullDiff))
	}

	for id, file := range orig {
		if file.Name() != fullDiff[id].Name() {
			t.Errorf("Diff should be identical as orig. Found diff in [%d]: %s vs %s",
				id, file.Name(), fullDiff[id].Name())
		}
	}
}
